/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/jo-pouradier/homelab-bot/master"
	"github.com/jo-pouradier/homelab-bot/utils"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list all connected agents.",
	Long:  `list all connected agents`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.SendCommandPipe("agent-ls", master.PipeName)
		fmt.Println(strings.ReplaceAll(utils.ReadCommandPipe(master.PipeName), "$", "\n"))
	},
}

func init() {
	agentCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
