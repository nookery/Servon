package caddy

import (
	"servon/core"

	"github.com/spf13/cobra"
)

type CommandOptions = core.CommandOptions

func Setup(app *core.App) {
	caddy := Caddy{
		App:     app,
		BaseDir: app.GetSoftwareRootFolder("caddy"),
	}

	app.RegisterGateway("caddy", &caddy)
	app.AddCommand(caddy.NewCaddyCommand(app))
}

func (c *Caddy) NewCaddyCommand(app *core.App) *cobra.Command {
	rootCmd := app.NewCommand(core.CommandOptions{
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
