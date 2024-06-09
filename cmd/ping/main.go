package main

import (
	"flag"

	"github.com/jo-pouradier/homelab-bot/agent"
	"golang.org/x/oauth2"
)

var (
	Address            = flag.String("addr", "localhost:50051", "the address to connect to like hostname:port")
	Tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	CaFile             = flag.String("ca_file", "./x509/ca_cert.pem", "The file containing the CA root cert file")
	ServerHostOverride = flag.String("server_host_override", "test.tbot.jo-pouradier.fr", "The server name used to verify the hostname returned by the TLS handshake")
)

func main() {
	flag.Parse()

	params := agent.NewAgentParams{
		Addr:               *Address,
		Tls:                *Tls,
		CaFile:             *CaFile,
		Token:              "some-secret-token",
		ServerHostOverride: *ServerHostOverride,
	}

	_, err := agent.NewAgent(params)

	if err != nil {
		return
	}
}
