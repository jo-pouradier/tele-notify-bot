package agent

import (
	"context"
	"fmt"
	"log"

	pb "github.com/jo-pouradier/homelab-bot/grpc"
	"github.com/jo-pouradier/homelab-bot/metrics"
)

type MetricsServerImpl struct {
	pb.UnimplementedMetricsServiceServer
}

func (s *MetricsServerImpl) Metrics(ctx context.Context, in *pb.Empty) (*pb.MetricsAllResponse, error) {
	log.Printf("Received: %v", in)
	cpu, _ := metrics.GetCPU1()
	mem, _ := metrics.GetMEM1()

	message := fmt.Sprintf("CPU: %.2f \nMEM: %.2f", cpu, mem)
	return &pb.MetricsAllResponse{Metrics: message}, nil
}
