package software

import (
	"github.com/spf13/cobra"
	"servon/core"
	"servon/core/utils"
)

// newListCmd 返回 list 子命令
func newListCmd(core *core.Core) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "显示支持的软件列表",
		Run: func(cmd *cobra.Command, args []string) {
			names := core.GetAllSoftware()

			utils.PrintList(names, "支持的软件列表")
		},
	}
}
