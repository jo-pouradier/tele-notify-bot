/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jo-pouradier/homelab-bot/master"
	"github.com/spf13/cobra"
)

// runMasterCmd represents the run command
var runMasterCmd = &cobra.Command{
	Use:   "run",
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

		// Handle graceful shutdown
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigs
			log.Println("Shutting down server...")
			a.DeleteNamedPipe()
			os.Exit(0)
		}()

		var wg sync.WaitGroup
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			if err := a.Serve(wg); err == nil {
				defer func() {
					if r := recover(); r != nil {
						a.DeleteNamedPipe()
					}
				}()
			}
		}(&wg)
		wg.Wait()

		a.ReadCmd()
	},
}

func init() {
	masterCmd.AddCommand(runMasterCmd)

	runMasterCmd.Flags().Bool("tls", false, "Use the TLS protocol")
	runMasterCmd.Flags().Int("port", 50051, "The server port")
	runMasterCmd.Flags().String("cert_file", "", "The TLS cert file path")
	runMasterCmd.Flags().String("key_file", "", "The TLS key file path")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
