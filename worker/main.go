package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os/exec"

	pb "thesis/proto"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The gRPC server port")
)

type server struct {
	pb.UnimplementedWorkerServer
}

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

func (s *server) RunJob(ctx context.Context, in *pb.JobParameters) (*pb.JobOutput, error) {
	param := in.GetMessage()
	output := runJob(param)
	return &pb.JobOutput{Output: output}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterWorkerServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
