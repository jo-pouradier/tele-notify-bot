package main

import (
	"github.com/jo-pouradier/homelab-bot/master"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()

}

func main() {
	master := master.NewMaster(true)
	master.Serve()
}
