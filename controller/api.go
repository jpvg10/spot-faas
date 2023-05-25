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

const (
	grpcPort = "50051"
)

func runJobInWorker(job Job) {
	var index int
	mu.Lock()
	for i := range jobs {
		if jobs[i].Id == job.Id {
			index = i
			break
		}
	}
	mu.Unlock()

	spotName := "spot" + job.Id
	var ip string

	if *local {
		ip = "localhost"
	} else {
		log.Printf("Creating spot VM: %v", spotName)
		ip = createVM(spotName)
		log.Printf("Spot VM %v created. The IP is: %v", spotName, ip)
	}

	ip = ip + ":" + grpcPort

	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewWorkerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	mu.Lock()
	jobs[index].Status = InProgress
	mu.Unlock()

	log.Printf("Launching job on spot VM")

	r, err := client.RunJob(ctx, &pb.JobParameters{Id: job.Id, Message: job.Message})
	if err != nil {
		log.Fatalf("Failed to run job: %v", err)
	}

	log.Printf("Completed job: %v\n", r.GetOutput())

	mu.Lock()
	jobs[index].Status = Completed
	jobs[index].Output = r.GetOutput()
	mu.Unlock()

	if !*local {
		log.Printf("Deleting spot VM: %v", spotName)
		deleteVM(spotName)
		log.Printf("Deleted spot VM: %v", spotName)
	}
}

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

	go runJobInWorker(newJob)

	c.IndentedJSON(http.StatusOK, newJob)
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

func pingApi(c *gin.Context) {
	c.Data(http.StatusOK, gin.MIMEHTML, []byte("API running.\n"))
}
