package software

import (
	"servon/core"
	"servon/core/utils"

	"github.com/spf13/cobra"
)

// Setup 注册到内核
func Setup(core *core.Core) error {
	core.AddCommand(GetSoftwareCommand(core))
	return nil
}

// GetSoftwareCommand 返回 software 命令
func GetSoftwareCommand(core *core.Core) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "software",
		Short: "软件管理",
		Long:  `📦 软件管理`,
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.PrintCommandHelp(cmd)
			return nil
		},
	}

	cmd.AddCommand(newListCmd(core))
	cmd.AddCommand(newInstallCmd(core))
	cmd.AddCommand(newInfoCmd(core))
	cmd.AddCommand(newStartCmd(core))
	cmd.AddCommand(newStopCmd(core))

	return cmd
}
