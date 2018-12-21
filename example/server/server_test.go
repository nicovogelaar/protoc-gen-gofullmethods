package main

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/myles-mcdonnell/blondie"
	"github.com/nicovogelaar/protoc-gen-gofullmethods/example"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	greeterClient example.GreeterClient
)

func TestMain(m *testing.M) {
	go start()

	opts := blondie.DefaultOptions()
	opts.QuietMode = true

	blondie.WaitForDeps([]blondie.DepCheck{blondie.NewTcpCheck("localhost", 8080, 30*time.Second)}, opts)

	greeterConn, err := grpc.Dial(":8080", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect to greeter server: %v", err)
	}
	defer func() {
		if err = greeterConn.Close(); err != nil {
			log.Printf("error while closing greeter server client connection: %v", err)
		}
	}()
	greeterClient = example.NewGreeterClient(greeterConn)

	os.Exit(m.Run())
}

func TestGreeterServer_SayHello(t *testing.T) {

	reply, err := greeterClient.SayHello(context.Background(), &example.HelloRequest{Name: "Nico"})
	if err != nil {
		t.Fatalf("failed to say hello: %v", err)
	}
	if reply.Message != "Hello Nico!" {
		t.Fatalf("unexpected reploy, got: %v", reply.Message)
	}
}

func TestGreeterServer_SayBye(t *testing.T) {
	_, err := greeterClient.SayBye(context.Background(), &example.ByeRequest{Name: "Nico"})
	if err == nil {
		t.Fatalf("unexpected reply. expected error")
	}
	s, ok := status.FromError(err)
	if !ok {
		t.Fatalf("unexpected error, got: %v", err)
	}
	if s.Code() != codes.PermissionDenied || s.Message() != "not allowed to say bye" {
		t.Fatalf("unexpected error, got: %v", err)
	}
}
