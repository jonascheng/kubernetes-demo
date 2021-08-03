package main

import (
	"context"
	"log"
	"net"
	"os"

	echo "github.com/jonascheng/kubernetes-demo/istio-grpc-lb/echo"
	"google.golang.org/grpc"
)

type server struct {
	echo.UnimplementedEchoServiceServer
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("listener fail: %v", err)
	}

	echoServer := grpc.NewServer()
	echo.RegisterEchoServiceServer(echoServer, &server{})
	if err := echoServer.Serve(listen); err != nil {
		log.Fatalf("server fail: %v", err)
	}
}

func (s *server) Echo(ctx context.Context, in *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{ServerAddress: os.Getenv("POD_IP")}, nil
}
