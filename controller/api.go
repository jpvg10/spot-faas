package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	pb "thesis/proto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	ip := "localhost"
	if !*local {
		ip = createVM("spot")
	}

	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewWorkerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.RunJob(ctx, &pb.JobParameters{Message: newJob.Message})
	if err != nil {
		log.Fatalf("Failed to get params: %v", err)
	}

	c.IndentedJSON(http.StatusOK, r.Output)
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
