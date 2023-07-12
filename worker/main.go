package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	pb "thesis/proto"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	port  = flag.Int("port", 50051, "The gRPC server port")
	image = flag.String("image", "jpvalencia/sum-fn", "The Docker image to run")
)

type server struct {
	pb.UnimplementedWorkerServiceServer
}

var sigChan = make(chan os.Signal, 1)
var interrupted = false

func runJob(args string, resultChan chan string, errorChan chan string) {
	dockerCommand := []string{"run"}

	if len(args) > 0 {
		dockerCommand = append(dockerCommand, "-e", fmt.Sprintf("FN_ARGS=%v", args))
	}
	dockerCommand = append(dockerCommand, "-e", "TIMEOUT=120000") // Function runs for 2 minutes
	dockerCommand = append(dockerCommand, *image)

	cmd := exec.Command("docker", dockerCommand...)

	var cmdOut bytes.Buffer
	cmd.Stdout = &cmdOut

	err := cmd.Run()
	if err != nil {
		errorChan <- err.Error()
		return
	}

	log.Printf("Container output: %s", cmdOut.String())
	resultChan <- cmdOut.String()
}

func (s *server) RunJob(ctx context.Context, in *pb.RunJobRequest) (*pb.RunJobResponse, error) {
	log.Printf("Received job request")

	if interrupted {
		// Interruption signal was received before
		log.Printf("Worker interrupted!")
		return &pb.RunJobResponse{Error: "Worker interrupted", Status: "failed"}, nil
	}

	args := in.GetArguments()

	resultChan := make(chan string, 1)
	errorChan := make(chan string, 1)
	go runJob(args, resultChan, errorChan)

	select {
	case result := <-resultChan:
		log.Printf("Job completed")
		return &pb.RunJobResponse{Result: result, Status: "completed"}, nil
	case err := <-errorChan:
		log.Printf("Job execution failed")
		log.Print(err)
		return &pb.RunJobResponse{Error: err, Status: "failed"}, nil
	case <-sigChan:
		log.Printf("Worker interrupted!")
		interrupted = true
		return &pb.RunJobResponse{Error: "Worker interrupted", Status: "failed"}, nil
	}
}

func (s *server) Ping(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()

	signal.Notify(sigChan, syscall.SIGTERM)

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
