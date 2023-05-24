package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var jobs = []Job{}
var mu sync.Mutex

func postMessage(c *gin.Context) {
	var newMessage Payload

	if err := c.BindJSON(&newMessage); err != nil {
		return
	}

	id := uuid.New()
	newJob := Job{Id: id.String(), Message: newMessage.Message, Status: Pending}

	mu.Lock()
	jobs = append(jobs, newJob)
	mu.Unlock()

	if !*local {
		go createVM("spot")
	}

	c.IndentedJSON(http.StatusCreated, newJob)
}

func getMessage(c *gin.Context) {
	id := c.Param("id")

	mu.Lock()
	defer mu.Unlock()

	for _, j := range jobs {
		if j.Id == id {
			c.IndentedJSON(http.StatusOK, j)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not found"})
}

func runApi() {
	router := gin.Default()
	router.GET("/message/:id", getMessage)
	router.POST("/message", postMessage)

	webAddress := "localhost:" + *webPort

	router.Run(webAddress)
}
