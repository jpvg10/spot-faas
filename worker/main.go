package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "thesis/poc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	controllerAddress = flag.String("controller", "localhost:50051", "The address to connect to")
)

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*controllerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewControllerClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetParams(ctx, &pb.PingMessage{Address: "worker"})
	if err != nil {
		log.Fatalf("Failed to get params: %v", err)
	}

	// Do the job
	jobOutput := runJob(r.GetName())

	// Send the output
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.SetOutput(ctx, &pb.JobOutput{Output: jobOutput})
	if err != nil {
		log.Fatal("Failed to send", err)
	}
}
