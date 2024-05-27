package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/jo-pouradier/homelab-bot/agent"
	pb "github.com/jo-pouradier/homelab-bot/grpc"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("Could not open port: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreetingServiceServer(s, &agent.PingServer{})
	pb.RegisterMetricsServiceServer(s, &agent.MetricsServerImpl{})

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
