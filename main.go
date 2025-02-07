package main

import (
	"os"
	"servon/cmd"
	"servon/cmd/software"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	// 导入插件
	_ "servon/plugins/caddy"
	_ "servon/plugins/clash"
	_ "servon/plugins/git"
	_ "servon/plugins/nodejs"
	_ "servon/plugins/npm"
	_ "servon/plugins/pnpm"
	_ "servon/plugins/yarn"
)

var (
	titleColor   = color.New(color.FgHiCyan, color.Bold)
	commandColor = color.New(color.FgHiGreen)
	infoColor    = color.New(color.FgHiWhite)
)

// RootCmd 是应用程序的根命令
var RootCmd = &cobra.Command{
	Use:   "servon",
	Short: "Servon - A lightweight server management panel",
	Long: titleColor.Sprintf("Servon - A lightweight server management panel\n\n") +
		infoColor.Sprintf("A comprehensive server management panel that provides\nboth CLI and Web interface for managing your server."),
}

func init() {
	// 加载所有插件
	if err := software.LoadPlugins(); err != nil {
		color.Red("Error loading plugins: %v\n", err)
		os.Exit(1)
	}

	// 注册所有子命令
	RootCmd.AddCommand(cmd.ServeCmd)
	RootCmd.AddCommand(cmd.VersionCmd)
	RootCmd.AddCommand(cmd.UpgradeCmd)
	RootCmd.AddCommand(cmd.DeployCmd)
	RootCmd.AddCommand(cmd.IPCmd)
	RootCmd.AddCommand(cmd.InstallCmd)
	RootCmd.AddCommand(cmd.SystemCmd)
	RootCmd.AddCommand(cmd.SoftwareCmd)

	// 添加 up 命令，并将 upgrade 设置为它的别名
	upCmd := *cmd.UpgradeCmd
	upCmd.Use = "up"
	upCmd.Aliases = []string{"upgrade"}
	RootCmd.AddCommand(&upCmd)

	// 设置彩色模板
	usageTemplate := titleColor.Sprintf("Usage:\n") +
		infoColor.Sprintf("  %s [command]\n", os.Args[0]) +
		titleColor.Sprintf("\nAvailable Commands:\n")

	// 添加命令列表
	RootCmd.Commands()
	for _, cmd := range RootCmd.Commands() {
		usageTemplate += commandColor.Sprintf("  %-15s", cmd.Name()) +
			infoColor.Sprintf("%s\n", cmd.Short)
	}

	usageTemplate += titleColor.Sprintf("\nFlags:\n") +
		infoColor.Sprintf("  -h, --help   help for %s\n\n", os.Args[0]) +
		infoColor.Sprintf("Use \"%s [command] --help\" for more information about a command.\n", os.Args[0])

	RootCmd.SetUsageTemplate(usageTemplate)
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		color.Red("Error: %v\n", err)
		os.Exit(1)
	}
}
