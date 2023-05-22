package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "thesis/poc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedControllerServer
}

func (s *server) GetParams(ctx context.Context, in *pb.PingMessage) (*pb.JobParameters, error) {
	p, _ := peer.FromContext(ctx)
	log.Printf("Received request from: %v", p.Addr)

	param := ""

	mu.Lock()
	defer mu.Unlock()

	if len(messages) > 0 {
		param = messages[0]
		messages = messages[1:]
		return &pb.JobParameters{Name: param}, nil
	} else {
		return nil, status.Error(codes.NotFound, "No params found")
	}
}

func (s *server) SetOutput(ctx context.Context, in *pb.JobOutput) (*emptypb.Empty, error) {
	log.Printf("Received output: %v", in.GetOutput())
	return &emptypb.Empty{}, nil
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
