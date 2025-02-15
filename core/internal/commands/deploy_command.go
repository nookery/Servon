package commands

import (
	"github.com/spf13/cobra"
)

// deployRootCommand 部署根命令
var deployRootCommand = NewCommand(CommandOptions{
	Use:     "deploy",
	Short:   "部署项目",
	Aliases: []string{"d"},
})

func GetDeployCommand() *cobra.Command {
	return deployRootCommand
}
