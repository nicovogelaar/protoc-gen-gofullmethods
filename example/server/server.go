package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/nicovogelaar/protoc-gen-gofullmethods/example"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type greeterServer struct {
}

func (server *greeterServer) SayHello(ctx context.Context, request *example.HelloRequest) (*example.HelloReply, error) {
	return &example.HelloReply{
		Message: fmt.Sprintf("Hello %v!", request.Name),
	}, nil
}

func (server *greeterServer) SayBye(ctx context.Context, request *example.ByeRequest) (*example.ByeReply, error) {
	return &example.ByeReply{
		Message: fmt.Sprintf("Bye %v!", request.Name),
	}, nil
}

func serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if info.FullMethod == example.Greeter_SayBye {
		return nil, status.Error(codes.PermissionDenied, "not allowed to say bye")
	}

	return handler(ctx, req)
}

func start() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(serverInterceptor))

	example.RegisterGreeterServer(server, &greeterServer{})
	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	start()
}
