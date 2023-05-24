package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "thesis/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	grpcPort = flag.Int("grpc-port", 50051, "The port for the gRPC server")
	webPort  = flag.String("web-port", "8080", "The port for the web server")
	local    = flag.Bool("local", true, "Use local or cloud worker")
)

type server struct {
	pb.UnimplementedControllerServer
}

func (s *server) GetParams(ctx context.Context, in *emptypb.Empty) (*pb.JobParameters, error) {
	p, _ := peer.FromContext(ctx)
	log.Printf("Received request from: %v", p.Addr)

	mu.Lock()
	defer mu.Unlock()

	for i, j := range jobs {
		if j.Status == Pending {
			jobs[i].Status = InProgress
			return &pb.JobParameters{Id: j.Id, Message: j.Message}, nil
		}
	}

	return nil, status.Error(codes.NotFound, "No params found")
}

func (s *server) SetOutput(ctx context.Context, in *pb.JobOutput) (*emptypb.Empty, error) {
	log.Printf("Received output: %v", in.GetOutput())

	mu.Lock()
	defer mu.Unlock()

	for i := range jobs {
		if jobs[i].Id == in.GetId() {
			jobs[i].Status = Completed
			jobs[i].Output = in.GetOutput()
			break
		}
	}

	if !*local {
		go deleteVM("spot")
	}

	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()

	go runApi()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
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
