package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type payload struct {
	Message string `json:"message"`
}

var messages = []string{}

func getMessages(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, messages)
}

func postMessage(c *gin.Context) {
	var newAlbum payload

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	messages = append(messages, newAlbum.Message)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func runApi() {
	router := gin.Default()
	router.GET("/messages", getMessages)
	// router.GET("/albums/:id", getAlbumByID)
	router.POST("/message", postMessage)

	router.Run("localhost:8080")
}
