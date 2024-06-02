package agent

import (
	"fmt"
	"log"
	"net"

	pb "github.com/jo-pouradier/homelab-bot/grpc"

	"google.golang.org/grpc"
)

type Agent interface {
	Serve()
}

type AgentImpl struct {
	lis net.Listener
	s   *grpc.Server
}

func NewAgent(port int) *AgentImpl {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatalf("Could not open port: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterGreetingServiceServer(s, &PingServerImpl{})
	pb.RegisterMetricsServiceServer(s, &MetricsServerImpl{})

	return &AgentImpl{
		lis: lis,
		s:   s,
	}

}

func (a *AgentImpl) Serve() {
	log.Printf("Server listening at %v", a.lis.Addr())
	if err := a.s.Serve(a.lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
