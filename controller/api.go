package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var jobs = []job{}
var mu sync.Mutex

func getMessages(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, jobs)
}

func postMessage(c *gin.Context) {
	var newMessage payload

	if err := c.BindJSON(&newMessage); err != nil {
		return
	}

	id := uuid.New()

	newJob := job{Id: id.String(), Message: newMessage.Message, Completed: false}

	mu.Lock()
	jobs = append(jobs, newJob)
	mu.Unlock()

	c.IndentedJSON(http.StatusCreated, newJob)
}

func getMessage(c *gin.Context) {
	id := c.Param("id")

	for _, a := range jobs {
		if a.Id == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
}

func runApi() {
	router := gin.Default()
	router.GET("/messages", getMessages)
	router.GET("/message/:id", getMessage)
	router.POST("/message", postMessage)

	router.Run("localhost:8080")
}
