package main

import (
	"context"
	"fmt"
	"hello/pkg/helloservice"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	Port = "0.0.0.0:3001"
)

type Server struct {
	helloservice.UnimplementedHelloServiceServer
}

func (s Server) Echo(ctx context.Context, req *helloservice.Request) (*helloservice.Response, error) {
	resp := &helloservice.Response{}
	resp.Message = req.Message

	log.Println("Got request...")

	time.Sleep(time.Millisecond * 100)
	return resp, nil
}

func main() {
	server := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     time.Second,
			MaxConnectionAge:      time.Second,
			MaxConnectionAgeGrace: time.Second,
			Time:                  time.Second * 10,
			Timeout:               time.Second,
		}),
	)
	defer server.GracefulStop()

	s := Server{}

	helloservice.RegisterHelloServiceServer(server, s)
	listener, err := net.Listen("tcp", Port)
	if err != nil {
		fmt.Println(err)
	}

	log.Println("Listening...")

	err = server.Serve(listener)
	if err != nil {
		fmt.Println(err)
	}
}
