package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/jo-pouradier/homelab-bot/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr               = flag.String("addr", "localhost:50051", "the address to connect to like hostname:port")
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverHostOverride = flag.String("server_host_override", "test.tbot.jo-pouradier.fr", "The server name used to verify the hostname returned by the TLS handshake")
)

func main() {
	flag.Parse()

	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = "./x509/ca_cert.pem"
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		log.Print("WARNING using insecure connection")
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(*addr, opts...)

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
