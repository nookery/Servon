package libs

import (
	"github.com/spf13/cobra"
)

type DeployManager struct{}

// deployRootCommand 部署根命令
var deployRootCommand = NewCommand(CommandOptions{
	Use:     "deploy",
	Short:   "部署项目",
	Aliases: []string{"d"},
})

func NewDeployManager() *DeployManager {
	return &DeployManager{}
}

func (c *DeployManager) GetDeployCommand() *cobra.Command {
	return deployRootCommand
}

// AppendDeploySubCommand 添加子命令
func (c *DeployManager) AppendDeploySubCommand(cmd *cobra.Command) {
	deployRootCommand.AddCommand(cmd)
}
