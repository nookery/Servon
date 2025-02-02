package main

import (
	"os"
	"servon/cmd"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// RootCmd 是应用程序的根命令
var RootCmd = &cobra.Command{
	Use:   "servon",
	Short: "Servon - A lightweight server management panel",
	Long: `Servon is a comprehensive server management panel that provides
both CLI and Web interface for managing your server.`,
}

func init() {
	// 注册所有子命令
	cmd.RegisterCommands(RootCmd)
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		color.Red("Error: %v\n", err)
		os.Exit(1)
	}
}
