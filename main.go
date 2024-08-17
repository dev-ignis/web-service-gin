package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"mht-web-service/handlers"
	"mht-web-service/utils"
)

func main() {
	// Set up S3 bucket and OpenAI API key from environment variables
	bucketName := os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		log.Fatalf("S3_BUCKET_NAME environment variable is not set")
	}

	openAiApiKey := os.Getenv("OPENAI_API_KEY")
	if openAiApiKey == "" {
		log.Fatalf("OPENAI_API_KEY environment variable is not set")
	}

	// Load AWS configuration
	cfg, s3Client := utils.LoadAWSConfig()
	fmt.Print(cfg)

	// Initialize the router
	router := gin.Default()

	// Set up routes
	handlers.RegisterAlbumRoutes(router, s3Client, bucketName)
	handlers.RegisterChatGptRoutes(router, openAiApiKey)

	// Run the server
	router.Run("0.0.0.0:8080")
}
