package cmd

import (
	"servon/cmd/software"

	"github.com/spf13/cobra"
)

var InstallCmd = &cobra.Command{
	Use:   "install [software]",
	Short: "Install software (alias for 'software install')",
	Long: `Install software is an alias for 'software install' command.
It provides a shorter way to install supported software.

Supported software:
  - caddy
  - node
  - pnpm
  - npm
  - clash`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 直接调用 software install 命令的执行函数
		return software.InstallCmd.RunE(cmd, args)
	},
	ValidArgsFunction: software.InstallCmd.ValidArgsFunction,
}

func init() {
	// 复制 software install 命令的所有标志
	InstallCmd.Flags().AddFlagSet(software.InstallCmd.Flags())
}
