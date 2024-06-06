package agent

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"

	pb "github.com/jo-pouradier/homelab-bot/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
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
		certs, err := tls.LoadX509KeyPair(params.CertFile, params.KeyFile)
		if err != nil {
			log.Fatalf("failed to load key pair: %s", err)
		}
		// creds, err := credentials.NewServerTLSFromFile(params.CertFile, params.KeyFile)
		// if err != nil {
		// 	log.Fatalf("Failed to generate credentials: %v", err)
		// }
		creds := credentials.NewServerTLSFromCert(&certs)

		opts = []grpc.ServerOption{
			grpc.UnaryInterceptor(ensureValidToken),
			grpc.Creds(creds),
		}
	} else {
		log.Print("WARNING your are not using tls Encryption, don't send any sensitive data")
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

// valid validates the authorization.
func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string.
	return token == "some-secret-token"
}

// ensureValidToken ensures a valid token exists within a request's metadata. If
// the token is missing or invalid, the interceptor blocks execution of the
// handler and returns an error. Otherwise, the interceptor invokes the unary
// handler.
func ensureValidToken(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	log.Printf("Metadata: %+v", md)
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}
