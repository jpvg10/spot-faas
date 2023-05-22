package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"time"

	pb "thesis/poc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	controllerAddress = flag.String("controller", "localhost:50051", "The address to connect to")
)

func runJob(param string) string {
	dockerCommand := []string{"run"}

	if len(param) > 0 {
		dockerCommand = append(dockerCommand, "-e", fmt.Sprintf("MESSAGE=%v", param))
	}
	dockerCommand = append(dockerCommand, "worker")

	cmd := exec.Command("docker", dockerCommand...)

	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	err := cmd.Run()

	if err != nil {
		log.Print(cmdErr.String())
		log.Fatal(err)
	}

	log.Printf("Container output: %s\n", cmdOut.String())
	return cmdOut.String()
}

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

	r, err := c.GetParams(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Failed to get params: %v", err)
	}

	// Do the job
	jobOutput := runJob(r.GetMessage())

	// Send the output
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.SetOutput(ctx, &pb.JobOutput{Id: r.GetId(), Output: jobOutput})
	if err != nil {
		log.Fatal("Failed to send", err)
	}
}
