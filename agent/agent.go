package agent

import (
	"context"
	"log"

	"github.com/google/uuid"
	pb "github.com/jo-pouradier/homelab-bot/grpc"
	"github.com/jo-pouradier/homelab-bot/metrics"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
)

type Agent interface {
	Serve()
	Ping(string)
	StreamMetrics()
}

type AgentImpl struct {
	ctx           context.Context
	conn          *grpc.ClientConn
	PingClient    *pb.GreetingServiceClient
	MetricsClient *pb.MetricsServiceClient
}

type NewAgentParams struct {
	Addr               string
	Tls                bool
	CaFile             string
	Token              string
	ServerHostOverride string
	AgentName          string
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

	ctx := context.WithoutCancel(context.Background())
	agentName := params.AgentName
	if agentName == "" {
		agentName = uuid.New().String()
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "name", agentName)

	return AgentImpl{
		ctx:  ctx,
		conn: conn,
	}, nil

}

func (agent *AgentImpl) Ping(msg string) {
	pingClient := pb.NewGreetingServiceClient(agent.conn)
	res, err := pingClient.Ping(agent.ctx, &pb.PingRequest{Name: msg})
	if err != nil {
		log.Fatalf("Error with rpc request: %v", err)
	}
	log.Printf("ping with txt=%v, response: %v", msg, res)
}

// TODO:
// add a retry strategy
// run inside a goroutine
// add mecanism to follow current status
func (agent *AgentImpl) StreamMetrics() {
	ctx := metadata.AppendToOutgoingContext(agent.ctx, "name", "metadata_name_testing", "data", "metrics")
	metricsClient := pb.NewMetricsServiceClient(agent.conn)

	for {
		stream, _ := metricsClient.GetMetricsStream(ctx, grpc.EmptyCallOption{})
		cpu, _ := metrics.GetCPU1()
		mem, _ := metrics.GetMEM1()
		log.Printf("New data cpu: %.2f, mem: %.2f", cpu, mem)
		if err := stream.Send(&pb.MetricsData{CpuPercentUsage: float32(cpu), MemPercentUsage: float32(mem)}); err != nil {
			log.Fatalf("error sending data: cpu: %.2f, mem: %.2f", cpu, mem)
		}
		res, _ := stream.Recv()
		if !res.AskMetrics {
			break
		}
	}
}
