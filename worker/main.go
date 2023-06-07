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
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	port  = flag.Int("port", 50051, "The gRPC server port")
	image = flag.String("image", "jpvalencia/worker", "The Docker image to run")
)

type server struct {
	pb.UnimplementedWorkerServiceServer
}

func runJob(args string) string {
	dockerCommand := []string{"run"}

	if len(args) > 0 {
		dockerCommand = append(dockerCommand, "-e", fmt.Sprintf("MESSAGE=%v", args))
	}
	dockerCommand = append(dockerCommand, *image)

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

func (s *server) RunJob(ctx context.Context, in *pb.RunJobRequest) (*pb.RunJobResponse, error) {
	args := in.GetArguments()
	result := runJob(args)
	return &pb.RunJobResponse{Result: result}, nil
}

func (s *server) Ping(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterWorkerServiceServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
