package managers

import (
	"github.com/spf13/cobra"
)

type DeployManager struct {
	deployRootCommand *cobra.Command
}

func NewDeployManager(rootCommand *cobra.Command) *DeployManager {
	return &DeployManager{
		deployRootCommand: rootCommand,
	}
}

// AppendDeploySubCommand 添加子命令
func (m *DeployManager) AppendDeploySubCommand(cmd *cobra.Command) {
	m.deployRootCommand.AddCommand(cmd)
}
