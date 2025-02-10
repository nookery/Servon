package software

import (
	"github.com/spf13/cobra"
	"servon/core"
)

// newListCmd 返回 list 子命令
func (p *SoftWarePlugin) newListCmd() *cobra.Command {
	return p.core.NewCommand(core.CommandOptions{
		Use:   "list",
		Short: "显示支持的软件列表",
		Run: func(cmd *cobra.Command, args []string) {
			names := p.core.GetAllSoftware()

			p.core.PrintList(names, "支持的软件列表")
		},
	})
}
