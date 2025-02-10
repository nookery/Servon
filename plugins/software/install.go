package software

import (
	"servon/core"

	"github.com/spf13/cobra"
)

func (p *SoftWarePlugin) newInstallCmd() *cobra.Command {
	cmd := p.core.NewCommand(core.CommandOptions{
		Use:   "install [software-name]",
		Short: "安装指定的软件",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			p.core.Install(args[0], nil)
		},
	})

	return cmd
}
