package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albumStore manages a list of albums.
type albumStore struct {
	mu     sync.RWMutex
	albums []album
}

// initialize the store with some data
func newAlbumStore() *albumStore {
	return &albumStore{
		albums: []album{
			{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
			{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
			{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
		},
	}
}

// getAlbums responds with the list of all albums as JSON.
func (s *albumStore) getAlbums(c *gin.Context) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c.IndentedJSON(http.StatusOK, s.albums)
}

// postAlbums adds an album from JSON received in the request body.
func (s *albumStore) postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	s.mu.Lock()
	s.albums = append(s.albums, newAlbum)
	s.mu.Unlock()

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id parameter sent by the client.
func (s *albumStore) getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, a := range s.albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func main() {
	store := newAlbumStore()

	router := gin.Default()
	router.GET("/albums", store.getAlbums)
	router.GET("/albums/:id", store.getAlbumByID)
	router.POST("/albums", store.postAlbums)

	router.Run("localhost:8080")
}
