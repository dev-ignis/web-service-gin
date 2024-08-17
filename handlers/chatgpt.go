package handlers

import (
	"example/web-service-gin/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatGptRequest struct {
	Prompt string `json:"message"`
}

func RegisterChatGptRoutes(router *gin.Engine, openAiApiKey string) {
	router.POST("/chatgpt", func(c *gin.Context) { handleChatGpt(c, openAiApiKey) })
}

func handleChatGpt(c *gin.Context, openAiApiKey string) {
	var request ChatGptRequest

	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	responseText, err := services.CallChatGptAPI(request.Prompt, openAiApiKey)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"response": responseText})
}
