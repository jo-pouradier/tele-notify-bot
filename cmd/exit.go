/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/jo-pouradier/homelab-bot/master"
	"github.com/jo-pouradier/homelab-bot/utils"
	"github.com/spf13/cobra"
)

// exitCmd represents the exit command
var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Kill the running server",
	Long:  `Kill the running server`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.SendCommandPipe("exit", master.PipeName)
	},
}

func init() {
	masterCmd.AddCommand(exitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
