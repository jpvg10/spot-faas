// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.15.8
// source: worker-service.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// WorkerServiceClient is the client API for WorkerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WorkerServiceClient interface {
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	RunJob(ctx context.Context, in *RunJobRequest, opts ...grpc.CallOption) (*RunJobResponse, error)
}

type workerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWorkerServiceClient(cc grpc.ClientConnInterface) WorkerServiceClient {
	return &workerServiceClient{cc}
}

func (c *workerServiceClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/worker.WorkerService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workerServiceClient) RunJob(ctx context.Context, in *RunJobRequest, opts ...grpc.CallOption) (*RunJobResponse, error) {
	out := new(RunJobResponse)
	err := c.cc.Invoke(ctx, "/worker.WorkerService/RunJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WorkerServiceServer is the server API for WorkerService service.
// All implementations must embed UnimplementedWorkerServiceServer
// for forward compatibility
type WorkerServiceServer interface {
	Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	RunJob(context.Context, *RunJobRequest) (*RunJobResponse, error)
	mustEmbedUnimplementedWorkerServiceServer()
}

// UnimplementedWorkerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedWorkerServiceServer struct {
}

func (UnimplementedWorkerServiceServer) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedWorkerServiceServer) RunJob(context.Context, *RunJobRequest) (*RunJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunJob not implemented")
}
func (UnimplementedWorkerServiceServer) mustEmbedUnimplementedWorkerServiceServer() {}

// UnsafeWorkerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WorkerServiceServer will
// result in compilation errors.
type UnsafeWorkerServiceServer interface {
	mustEmbedUnimplementedWorkerServiceServer()
}

func RegisterWorkerServiceServer(s grpc.ServiceRegistrar, srv WorkerServiceServer) {
	s.RegisterService(&WorkerService_ServiceDesc, srv)
}

func _WorkerService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkerServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/worker.WorkerService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkerServiceServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkerService_RunJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkerServiceServer).RunJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/worker.WorkerService/RunJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkerServiceServer).RunJob(ctx, req.(*RunJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// WorkerService_ServiceDesc is the grpc.ServiceDesc for WorkerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WorkerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "worker.WorkerService",
	HandlerType: (*WorkerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _WorkerService_Ping_Handler,
		},
		{
			MethodName: "RunJob",
			Handler:    _WorkerService_RunJob_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "worker-service.proto",
}
