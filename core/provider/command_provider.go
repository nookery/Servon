package provider

import (
	"servon/core/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	titleColor = color.New(color.FgHiCyan, color.Bold)
	infoColor  = color.New(color.FgHiWhite)
)

// CommandProvider 命令行命令执行器
type CommandProvider struct {
	RootCmd *cobra.Command
}

// NewCommandProvider 创建命令行命令执行器
func NewCommandProvider() CommandProvider {
	rootCmd := &cobra.Command{
		Use:   "servon",
		Short: "Servon - A lightweight server management panel",
		Long: titleColor.Sprintf(`
   ____                          
  / ___|  ___ _ ____   ___  _ __  
  \___ \ / _ \ '__\ \ / / || '_ \ 
   ___) |  __/ |   \ V /| || | | |
  |____/ \___|_|    \_/ |_||_| |_|
                                  
`) +
			titleColor.Sprintf("🚀 Servon - A lightweight server management panel\n") +
			titleColor.Sprintf("🌟 Servon - 轻量级服务器管理面板\n\n") +
			infoColor.Sprintf("📦 A comprehensive server management panel that provides\n   both CLI and Web interface for managing your server.\n") +
			infoColor.Sprintf("🔧 一个提供命令行和Web界面的全功能服务器管理面板。"),
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrintCommandHelp(cmd)
		},
	}

	return CommandProvider{
		RootCmd: rootCmd,
	}
}

// AddCommand 添加命令
func (p *CommandProvider) AddCommand(cmd *cobra.Command) {
	p.RootCmd.AddCommand(cmd)
}
