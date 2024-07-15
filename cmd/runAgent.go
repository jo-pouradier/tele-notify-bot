/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/jo-pouradier/homelab-bot/agent"
	"github.com/spf13/cobra"
)

// agentCmd represents the agent command
var runAgentCmd = &cobra.Command{
	Use:   "run",
	Short: "A bk sief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve the flags
		address, err := cmd.Flags().GetString("addr")
		if err != nil {
			log.Fatalf("Could not get 'addr' flag: %v", err)
		}

		tls, err := cmd.Flags().GetBool("tls")
		if err != nil {
			log.Fatalf("Could not get 'tls' flag: %v", err)
		}

		caFile, err := cmd.Flags().GetString("ca_file")
		if err != nil {
			log.Fatalf("Could not get 'ca_file' flag: %v", err)
		}

		serverHostOverride, err := cmd.Flags().GetString("server_host_override")
		if err != nil {
			log.Fatalf("Could not get 'server_host_override' flag: %v", err)
		}

		agentName, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatalf("Could not get 'name' flag: %v", err)
		}

		// Create the agent parameters
		params := agent.NewAgentParams{
			Addr:               address,
			Tls:                tls,
			CaFile:             caFile,
			Token:              "some-secret-token", // replace with your token management
			ServerHostOverride: serverHostOverride,
			AgentName:          agentName,
		}

		// Initialize the agent
		agent, err := agent.NewAgent(params)
		if err != nil {
			log.Fatal(err)
		}

		// Use the agent
		agent.Ping("test")
		agent.Ping("ping")
		agent.StreamMetrics()
	},
}

func init() {
	agentCmd.AddCommand(runAgentCmd)

	runAgentCmd.Flags().String("name", "", "Name of the agent for for better recognition, if empty it will be a UUID")
	runAgentCmd.Flags().String("addr", "localhost:50051", "The address to connect to like hostname:port")
	runAgentCmd.Flags().Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	runAgentCmd.Flags().String("ca_file", "./x509/ca_cert.pem", "The file containing the CA root cert file")
	runAgentCmd.Flags().String("server_host_override", "test.tbot.jo-pouradier.fr", "The server name used to verify the hostname returned by the TLS handshake")
}
