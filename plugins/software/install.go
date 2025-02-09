package software

import (
	"github.com/spf13/cobra"
	"servon/core"
	"servon/core/utils"
)

func newInstallCmd(core *core.Core) *cobra.Command {
	return &cobra.Command{
		Use:   "install [软件名称]",
		Short: "安装指定的软件",
		Long:  `安装指定的软件。`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				utils.PrintCommandHelp(cmd)
				return
			}

			name := args[0]
			err := core.InstallSoftware(name, nil)
			if err != nil {
				utils.PrintError(err)
			}
		},
	}
}
