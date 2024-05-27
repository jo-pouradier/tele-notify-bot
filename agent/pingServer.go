package agent

import (
	"context"
	"log"

	pb "github.com/jo-pouradier/homelab-bot/grpc"
)

type PingServer struct {
	pb.UnimplementedGreetingServiceServer
}

func (s *PingServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PongResponse, error) {
	log.Printf("Received: %v", in.GetName())
	message := "What did you send me ?"
	if in.GetName() == "ping" {
		message = "pong"
	}
	return &pb.PongResponse{Message: message}, nil
}
