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

const (
	defaultName = "localhost:50052"
)

var (
	addr = flag.String("addr", "localhost:50051", "The address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
	port = flag.Int("port", 50052, "The server port")
)

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewControllerClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetParams(ctx, &pb.PingMessage{Address: *name})
	if err != nil {
		log.Fatalf("Could not get params: %v", err)
	}

	// do the job
	jobOutput := runJob(r.GetName())

	// send output
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.SetOutput(ctx, &pb.JobOutput{Output: jobOutput})
	if err != nil {
		log.Fatal("Failed to send", err)
	}
}
