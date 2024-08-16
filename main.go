package main

import (
	"context"
	"encoding/json"
	"fmt"
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

var s3Client *s3.Client
var bucketName string

func main() {
	// Set the bucket name from an environment variable
	bucketName = os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		log.Fatalf("S3_BUCKET_NAME environment variable is not set")
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
