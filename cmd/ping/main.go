package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/jo-pouradier/homelab-bot/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Erro creating new Client: %v", err)
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

}
