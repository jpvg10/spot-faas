package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "thesis/poc/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedControllerServer
}

func (s *server) GetParams(ctx context.Context, in *pb.PingMessage) (*pb.JobParameters, error) {
	log.Printf("Sent params: %v", in.GetAddress())
	param := ""
	mu.Lock()
	if len(messages) > 0 {
		param = messages[0]
		messages = messages[1:]
	}
	mu.Unlock()
	return &pb.JobParameters{Name: param}, nil
}

func (s *server) SetOutput(ctx context.Context, in *pb.JobOutput) (*empty.Empty, error) {
	log.Printf("Received output: %v", in.GetOutput())
	return &empty.Empty{}, nil
}

func main() {
	go runApi()

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterControllerServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
