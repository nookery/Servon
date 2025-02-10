package caddy

import (
	"servon/core"

	"github.com/spf13/cobra"
)

type CommandOptions = core.CommandOptions

func NewCaddyCommand(core *core.Core) *cobra.Command {
	rootCmd := core.NewCommand(CommandOptions{
		Use:   "caddy",
		Short: "Caddy 管理命令",
	})

	caddy := Caddy{
		Core: core,
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
