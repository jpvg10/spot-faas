package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type payload struct {
	Message string `json:"message"`
}

var messages = []string{}
var mu sync.Mutex

func getMessages(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, messages)
}

func postMessage(c *gin.Context) {
	var newAlbum payload

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	mu.Lock()
	messages = append(messages, newAlbum.Message)
	mu.Unlock()
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func runApi() {
	router := gin.Default()
	router.GET("/messages", getMessages)
	// router.GET("/albums/:id", getAlbumByID)
	router.POST("/message", postMessage)

	router.Run("localhost:8080")
}
