package main

import (
	"flag"

	"github.com/jo-pouradier/homelab-bot/agent"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()

	a := agent.NewAgent(*port)
	a.Serve()
}
