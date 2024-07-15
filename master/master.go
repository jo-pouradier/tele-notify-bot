package master

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"syscall"

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

	PipeName = "/tmp/tele-notify-bot-serve"
)

type Master interface {
	Serve()
}

type MasterImpl struct {
	lis           net.Listener
	Servers       *grpc.Server
	MetricsServer *MetricsServerImpl
	PipeName      string
	mu            *sync.Mutex
}

type NewMasterParams struct {
	Port     int
	Tls      bool
	CertFile string
	KeyFile  string
}

func NewMaster(params NewMasterParams) *MasterImpl {
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
	metricsServer := &MetricsServerImpl{Agents: make(map[string]*pb.MetricsData), mu: &sync.Mutex{}}
	pb.RegisterMetricsServiceServer(s, metricsServer)

	return &MasterImpl{
		lis:           lis,
		Servers:       s,
		PipeName:      PipeName,
		MetricsServer: metricsServer,
		mu:            &sync.Mutex{},
	}

}

func (a *MasterImpl) Serve(wg *sync.WaitGroup) error {
	a.createNamedPipe()
	wg.Done()
	log.Printf("Server listening at %v", a.lis.Addr())
	if err := a.Servers.Serve(a.lis); err != nil {
		// log.Fatalf("failed to serve: %v", err)
		return errors.New(fmt.Sprintf("failed to serve: %v", err))
	}
	return nil
}

func (a *MasterImpl) createNamedPipe() {
	if err := syscall.Mkfifo(a.PipeName, 0660); err != nil {
		log.Fatalf("failed to create named pipe: %v", err)
	}
}

func (a *MasterImpl) DeleteNamedPipe() {
	if err := os.Remove(a.PipeName); err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to remove pipe %s: %v", a.PipeName, err)
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
