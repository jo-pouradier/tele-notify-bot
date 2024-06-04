package main

import (
	"flag"

	"github.com/jo-pouradier/homelab-bot/agent"
)

var (
	port     = flag.Int("port", 50051, "The server port")
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file: -certFile=... ")
	keyFile  = flag.String("key_file", "", "The TLS key file")
)

func main() {
	flag.Parse()

	a := agent.NewAgent(agent.NewAgentParams{
		Port:     *port,
		Tls:      *tls,
		CertFile: *certFile,
		KeyFile:  *keyFile,
	})
	a.Serve()
}
