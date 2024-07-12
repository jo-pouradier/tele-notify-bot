/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/jo-pouradier/homelab-bot/master"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve the flags
		tls, err := cmd.Flags().GetBool("tls")
		if err != nil {
			log.Fatalf("Could not get 'tls' flag: %v", err)
		}

		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			log.Fatalf("Could not get 'port' flag: %v", err)
		}

		certFile, err := cmd.Flags().GetString("cert_file")
		if err != nil {
			log.Fatalf("Could not get 'cert_file' flag: %v", err)
		}

		keyFile, err := cmd.Flags().GetString("key_file")
		if err != nil {
			log.Fatalf("Could not get 'key_file' flag: %v", err)
		}

		// create server and serve
		a := master.NewMaster(master.NewMasterParams{
			Port:     port,
			Tls:      tls,
			CertFile: certFile,
			KeyFile:  keyFile,
		})
		a.Serve()

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().Bool("tls", false, "Use the TLS protocol")
	serveCmd.Flags().Int("port", 50051, "The server port")
	serveCmd.Flags().String("cert_file", "", "The TLS cert file path")
	serveCmd.Flags().String("key_file", "", "The TLS key file path")

	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
