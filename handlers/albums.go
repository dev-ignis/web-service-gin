package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func RegisterAlbumRoutes(router *gin.Engine, s3Client *s3.Client, bucketName string) {
	router.GET("/albums/:id", func(c *gin.Context) { getAlbumByID(c, s3Client, bucketName) })
	router.POST("/albums", func(c *gin.Context) { postAlbums(c, s3Client, bucketName) })
}

func postAlbums(c *gin.Context, s3Client *s3.Client, bucketName string) {
	var newAlbum Album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	albumJSON, err := json.Marshal(newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	reader := strings.NewReader(string(albumJSON))
	objectKey := fmt.Sprintf("albums/%s.json", newAlbum.ID)
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   reader,
	})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumByID(c *gin.Context, s3Client *s3.Client, bucketName string) {
	id := c.Param("id")
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

	var album Album
	if err := json.NewDecoder(resp.Body).Decode(&album); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, album)
}
