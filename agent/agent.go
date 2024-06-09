package agent

import (
	"context"
	"log"
	"time"

	pb "github.com/jo-pouradier/homelab-bot/grpc"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/oauth"
)

type AgentImpl struct {
	conn *grpc.ClientConn
}

type NewAgentParams struct {
	Addr               string
	Tls                bool
	CaFile             string
	Token              string
	ServerHostOverride string
}

func NewAgent(params NewAgentParams) (AgentImpl, error) {
	var opts []grpc.DialOption
	if params.Tls {
		if params.CaFile == "" {
			params.CaFile = "./x509/ca_cert.pem"
		}
		creds, err := credentials.NewClientTLSFromFile(params.CaFile, params.ServerHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
		// authentication
		perRPC := oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: params.Token})}
		opts = append(opts, grpc.WithPerRPCCredentials(perRPC))
	} else {
		log.Print("WARNING using insecure connection")
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(params.Addr, opts...)

	if err != nil {
		log.Fatalf("Error connectiong to server: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreetingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.Ping(ctx, &pb.PingRequest{Name: "ping"})
	if err != nil {
		log.Fatalf("Error with rpc request: %v", err)
	}
	log.Printf("ping 1 with txt=ping: %v", res)

	res2, _ := c.Ping(ctx, &pb.PingRequest{Name: "test"})
	log.Printf("ping 2 with txt=test: %v", res2)

	m := pb.NewMetricsServiceClient(conn)
	metricsCtx, metricsCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer metricsCancel()
	metrics, _ := m.Metrics(metricsCtx, &pb.Empty{})
	log.Printf("get metrics: %s", metrics)

	return AgentImpl{conn: conn}, nil

}
