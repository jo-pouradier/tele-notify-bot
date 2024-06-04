package agent

import (
	"fmt"
	"log"
	"net"

	pb "github.com/jo-pouradier/homelab-bot/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Agent interface {
	Serve()
}

type AgentImpl struct {
	lis net.Listener
	s   *grpc.Server
}

type NewAgentParams struct {
	Port     int
	Tls      bool
	CertFile string
	KeyFile  string
}

func NewAgent(params NewAgentParams) *AgentImpl {
	if params.Port == 0 {
		params.Port = 50000
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", params.Port))

	if err != nil {
		log.Fatalf("Could not open port: %v", err)
	}

	var opts []grpc.ServerOption
	if params.Tls {
		if params.CertFile == "" {
			params.CertFile = "./x509/server_cert.pem"
		}
		if params.KeyFile == "" {
			params.KeyFile = "x509/server_key.pem"
		}
		creds, err := credentials.NewServerTLSFromFile(params.CertFile, params.KeyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials: %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	s := grpc.NewServer(opts...)
	// s := grpc.NewServer()

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
