package master

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	pb "github.com/jo-pouradier/homelab-bot/grpc"
	"github.com/jo-pouradier/homelab-bot/metrics"
	"google.golang.org/grpc/metadata"
)

type MetricsServerImpl struct {
	pb.UnimplementedMetricsServiceServer

	mu *sync.Mutex

	Agents map[string]*pb.MetricsData
}

func (s *MetricsServerImpl) Metrics(ctx context.Context, in *pb.Empty) (*pb.MetricsAllResponse, error) {
	log.Printf("Received: %v", in)
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Printf("metadata from client: %+v", md)
	}
	cpu, _ := metrics.GetCPU1()
	mem, _ := metrics.GetMEM1()

	message := fmt.Sprintf("CPU: %.2f \nMEM: %.2f", cpu, mem)
	return &pb.MetricsAllResponse{Metrics: message}, nil
}

func (s *MetricsServerImpl) GetMetricsStream(streamMetrics pb.MetricsService_GetMetricsStreamServer) error {
	ctx, cancel := context.WithCancel(streamMetrics.Context())
	defer cancel()

	// read metadata
	md, ok := metadata.FromIncomingContext(streamMetrics.Context())
	if !ok {
		return fmt.Errorf("missing metadata")
	}

	name := md.Get("name")[0]
	log.Printf("get name agent: %s", name)
	s.mu.Lock()
	s.Agents[name] = &pb.MetricsData{}
	s.mu.Unlock()

	go func() {
		<-ctx.Done()
		s.mu.Lock()
		delete(s.Agents, name)
		s.mu.Unlock()
		log.Printf("Agent disconnected: %s", name)
	}()

	for {
		in, err := streamMetrics.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			cancel()
			return err
		}

		s.mu.Lock()
		s.Agents[name] = in
		s.mu.Unlock()

		log.Printf("Data stream from %s: %+v", name, in)

		time.Sleep(5 * time.Second)
		if err := streamMetrics.Send(&pb.AskMetrics{AskMetrics: true}); err != nil {
			cancel()
			return err
		}
	}

}
