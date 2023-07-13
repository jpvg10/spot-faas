package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	pb "thesis/proto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var jobs = []Job{}
var mu sync.Mutex

const (
	grpcPort = "50051"
)

func setError(index int, message string) {
	mu.Lock()
	jobs[index].Status = "failed"
	jobs[index].Error = message
	mu.Unlock()
}

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

	spotName := "spot-" + job.Id
	var ip string

	if *local {
		ip = "localhost"
	} else {
		log.Printf("%v - Creating spot VM", spotName)
		ip, createErr := createVM(spotName)
		if createErr != nil {
			log.Printf("%v - Failed to create the spot VM: %v", spotName, createErr)
			setError(index, createErr.Error())
			return
		}
		log.Printf("%v - Spot VM created. The IP is: %v", spotName, ip)
	}

	ip = ip + ":" + grpcPort

	conn, dialErr := grpc.Dial(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if dialErr != nil {
		log.Printf("%v - Did not connect: %v", spotName, dialErr)
		setError(index, dialErr.Error())
		return
	}
	defer conn.Close()
	client := pb.NewWorkerServiceClient(conn)

	for i := 0; ; i++ {
		time.Sleep(time.Second)
		if i%5 == 0 {
			log.Printf("%v - Attempting to contact the worker: %v", spotName, i)
		}

		ctxPing, cancelPing := context.WithTimeout(context.Background(), time.Second)
		defer cancelPing()

		_, pingErr := client.Ping(ctxPing, &emptypb.Empty{})
		if pingErr == nil {
			break
		} else if i >= 60 {
			log.Printf("%v - Failed to contact the worker in 1 minute", spotName)
			setError(index, "Failed to contact the worker in 1 minute")
			return
		}
	}

	log.Printf("%v - Launching job on spot VM", spotName)

	mu.Lock()
	jobs[index].Status = InProgress
	mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()

	r, runErr := client.RunJob(ctx, &pb.RunJobRequest{Id: job.Id, Arguments: job.Arguments})
	if runErr != nil {
		log.Printf("%v - Job returned an error: %v", spotName, runErr)
		setError(index, runErr.Error())
		return
	}

	resultString := r.GetResult()
	status := r.GetStatus()
	log.Printf("%v - Job status: %v", spotName, status)
	log.Printf("%v - Job result: %v", spotName, resultString)

	var resultJson map[string]interface{}
	unmarshalErr := json.Unmarshal([]byte(resultString), &resultJson)

	mu.Lock()
	jobs[index].Status = StatusType(status)

	if unmarshalErr != nil {
		// Result string
		jobs[index].Result = r.GetResult()
	} else {
		// Result JSON
		jobs[index].Result = resultJson
	}
	jobs[index].Error = r.GetError()
	mu.Unlock()

	if !*local {
		log.Printf("%v - Deleting spot VM", spotName)
		deleteErr := deleteVM(spotName)
		if deleteErr != nil {
			log.Printf("%v - Failed to delete the spot VM: %v", spotName, deleteErr)
		} else {
			log.Printf("%v - Deleted spot VM", spotName)
		}
	}
}

func postJob(c *gin.Context) {
	var args string
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		args = ""
	} else {
		args = string(jsonData)
	}

	id := uuid.New()
	newJob := Job{Id: id.String(), Arguments: args, Status: Pending}

	mu.Lock()
	jobs = append(jobs, newJob)
	mu.Unlock()

	go runJobInWorker(newJob)

	c.IndentedJSON(http.StatusOK, newJob)
}

func getJob(c *gin.Context) {
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
