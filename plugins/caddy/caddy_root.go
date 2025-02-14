package caddy

import (
	"servon/core"

	"github.com/spf13/cobra"
)

type CommandOptions = core.CommandOptions

func Setup(core *core.Core) {
	caddy := Caddy{
		Core:    core,
		BaseDir: core.DataManager.GetSoftwareRootFolder("caddy"),
	}

	core.RegisterSoftware("caddy", &caddy)
	core.AddCommand(caddy.NewCaddyCommand(core))
}

func (c *Caddy) NewCaddyCommand(core *core.Core) *cobra.Command {
	rootCmd := core.NewCommand(CommandOptions{
		Use:   "caddy",
		Short: "Caddy 管理命令",
	})

	// Add subcommands
	rootCmd.AddCommand(
		c.NewInstallCommand(),
		c.NewProxyCommand(),
	)

	return rootCmd
}
