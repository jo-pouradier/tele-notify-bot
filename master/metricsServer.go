package master

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	pb "github.com/jo-pouradier/homelab-bot/grpc"
	"github.com/jo-pouradier/homelab-bot/logger"
	"github.com/jo-pouradier/homelab-bot/metrics"
)

type MetricsServerImpl struct {
	pb.UnimplementedMetricsServiceServer

	mu      sync.Mutex // protects routeNotes
	metrics *pb.MetricsData
}

func (s *MetricsServerImpl) Metrics(ctx context.Context, in *pb.Empty) (*pb.MetricsAllResponse, error) {
	log.Printf("Received: %v", in)
	cpu, _ := metrics.GetCPU1()
	mem, _ := metrics.GetMEM1()

	message := fmt.Sprintf("CPU: %.2f \nMEM: %.2f", cpu, mem)
	return &pb.MetricsAllResponse{Metrics: message}, nil
}

func (s *MetricsServerImpl) GetMetricsStream(streamMetrics pb.MetricsService_GetMetricsStreamServer) error {
	for {
		in, err := streamMetrics.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		s.mu.Lock()
		s.metrics = in
		s.mu.Unlock()

		// log.Printf("Data stream: %+v", in)
		logger.Debug("Data stream: %+v", in)

		time.Sleep(5 * time.Second)
		streamMetrics.Send(&pb.AskMetrics{AskMetrics: true})
	}

}
