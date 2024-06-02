package main

import (
	"flag"

	"github.com/jo-pouradier/homelab-bot/master"
	"github.com/joho/godotenv"
)

var debug = flag.Bool("debug", true, "Use debug Mode, default: -debug=true")

func init() {
	godotenv.Load()
	flag.Parse()
}

func main() {
	master := master.NewMaster(*debug)
	master.Serve()
}
