package main

import (
	"context"
	"fmt"
	"hello/pkg/helloservice"
	"net"
	"time"

	"google.golang.org/grpc"
)

const (
	Port = ":3001"
)

type Server struct {
	helloservice.UnimplementedHelloServiceServer
}

func (s Server) Echo(ctx context.Context, req *helloservice.Request) (*helloservice.Response, error) {
	resp := &helloservice.Response{}
	resp.Message = req.Message

	time.Sleep(time.Millisecond * 100)
	return resp, nil
}

func main() {
	server := grpc.NewServer()
	defer server.GracefulStop()

	s := Server{}

	helloservice.RegisterHelloServiceServer(server, s)
	listener, err := net.Listen("tcp", Port)
	if err != nil {
		fmt.Println(err)
	}

	err = server.Serve(listener)
	if err != nil {
		fmt.Println(err)
	}
}
