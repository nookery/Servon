package software

import (
	"fmt"

	"github.com/spf13/cobra"
	"servon/core"
)

func newInstallCmd(core *core.Core) *cobra.Command {
	return &cobra.Command{
		Use:   "install [软件名称]",
		Short: "安装指定的软件",
		Long:  `安装指定的软件。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("请指定要安装的软件名称")
			}

			name := args[0]
			return core.InstallSoftware(name, nil)
		},
	}
}
