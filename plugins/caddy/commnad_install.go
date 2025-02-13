package caddy

import (
	"github.com/spf13/cobra"
)

func NewInstallCommand(caddy *Caddy) *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "安装 Caddy",
		RunE: func(cmd *cobra.Command, args []string) error {
			return caddy.Install()
		},
	}
}
