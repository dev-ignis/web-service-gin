package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type chatGptRequest struct {
	Prompt string `json:"message"`
}

type chatGptResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string  `json:"role"`
			Content string  `json:"content"`
			Refusal *string `json:"refusal"`
		} `json:"message"`
		Logprobs     *interface{} `json:"logprobs"`
		FinishReason string       `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

var s3Client *s3.Client
var bucketName string
var openAiApiKey string

func main() {
	// Set the bucket name from an environment variable
	bucketName = os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		log.Fatalf("S3_BUCKET_NAME environment variable is not set")
	}

	// Set the OpenAI API key from an environment variable
	openAiApiKey = os.Getenv("OPENAI_API_KEY")
	if openAiApiKey == "" {
		log.Fatalf("OPENAI_API_KEY environment variable is not set")
	}

	// Load the AWS SDK configuration using environment variables
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create an S3 client
	s3Client = s3.NewFromConfig(cfg)

	// Initialize the router
	router := gin.Default()

	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.POST("/chatgpt", handleChatGpt)

	router.Run("0.0.0.0:8080")
}

func postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Convert album to JSON
	albumJSON, err := json.Marshal(newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a reader from the JSON data
	reader := strings.NewReader(string(albumJSON))

	// Upload the JSON data to S3
	objectKey := fmt.Sprintf("albums/%s.json", newAlbum.ID)
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   reader, // io.Reader is used here directly
	})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Get the album JSON from S3
	objectKey := fmt.Sprintf("albums/%s.json", id)
	resp, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}

	defer resp.Body.Close()

	var a album
	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, a)
}

func handleChatGpt(c *gin.Context) {
	var request chatGptRequest

	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call the OpenAI API with the provided prompt
	responseText, err := callChatGptAPI(request.Prompt)
	if err != nil {
		// Return the actual error from ChatGPT
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"response": responseText})
}

func callChatGptAPI(prompt string) (string, error) {
	apiURL := "https://api.openai.com/v1/chat/completions"

	reqBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens":  100,
		"temperature": 0.7,
	})
	if err != nil {
		log.Printf("Error marshaling request body: %v", err)
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openAiApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request to OpenAI: %v", err)
		return "", fmt.Errorf("error sending request to OpenAI: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-200 response from OpenAI API: %s - %s", resp.Status, string(body))
		return "", fmt.Errorf("received non-200 response: %s - %s", resp.Status, string(body))
	}

	var chatGptResp chatGptResponse
	err = json.Unmarshal(body, &chatGptResp)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	if len(chatGptResp.Choices) > 0 {
		return strings.TrimSpace(chatGptResp.Choices[0].Message.Content), nil
	}

	log.Printf("No response from ChatGPT")
	return "", fmt.Errorf("no response from ChatGPT")
}
