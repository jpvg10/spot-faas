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

var clients []string

func (s *server) GetParams(ctx context.Context, in *pb.PingMessage) (*pb.JobParameters, error) {
	log.Printf("Connection: %v", in.GetAddress())
	clients = append(clients, in.GetAddress())
	fmt.Println(clients)

	return &pb.JobParameters{Name: "job"}, nil
}

func (s *server) SetOutput(ctx context.Context, in *pb.JobOutput) (*empty.Empty, error) {
	log.Printf("Received output: %v", in.GetOutput())
	return &empty.Empty{}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterControllerServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
