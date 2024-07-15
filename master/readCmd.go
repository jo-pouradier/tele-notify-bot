package master

import (
	"log"
	"os"
	"strings"

	"github.com/jo-pouradier/homelab-bot/utils"
)

// goroutine to read the current master named pipe and respond
func (master *MasterImpl) ReadCmd() {
	// Read commands from named pipe
	for {
		command := utils.ReadCommandPipe(master.PipeName)
		log.Printf("Received command: %s", command)

		// Handle the command
		switch command {
		case "exit":
			master.Exit()
		case "agent-ls":
			allAgent := master.AgentLs()
			log.Println(allAgent)
			str := strings.Join(allAgent, "$")
			log.Println(str)
			utils.SendCommandPipe(str, master.PipeName)
		default:
			continue
		}
	}
}

func (master *MasterImpl) Exit() {
	log.Println("Exiting server...")
	master.DeleteNamedPipe()
	os.Exit(0)
}

// TODO: split func to only use one mu at a time
func (master *MasterImpl) AgentLs() []string {
	master.mu.Lock()
	master.MetricsServer.mu.Lock()
	defer master.MetricsServer.mu.Unlock()
	defer master.mu.Unlock()

	keys := make([]string, len(master.MetricsServer.Agents))
	i := 0
	for k := range master.MetricsServer.Agents {
		keys[i] = k
		i++
	}
	return keys
}
