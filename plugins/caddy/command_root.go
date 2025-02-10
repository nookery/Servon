package caddy

import (
	"servon/core"

	"github.com/spf13/cobra"
)

func NewCaddyCommand(core *core.Core) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "caddy",
		Short: "Caddy 管理命令",
		Long:  "Caddy 管理命令",
		RunE: func(cmd *cobra.Command, args []string) error {
			core.PrintCommandHelp(cmd)
			return nil
		},
	}

	caddy := Caddy{
		core: core,
		CaddyConfig: CaddyConfig{
			BaseDir: "/data/caddy",
		},
	}

	// Add subcommands
	rootCmd.AddCommand(
		NewInstallCommand(&caddy),
		NewProxyCommand(&caddy),
	)

	return rootCmd
}
