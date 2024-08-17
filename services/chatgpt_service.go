package services

import (
	"bytes"
	"encoding/json"
	"example/web-service-gin/models"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func CallChatGptAPI(prompt, openAiApiKey string) (string, error) {
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

	var chatGptResp models.ChatGptResponse
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
