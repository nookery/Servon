package providers

import (
	"servon/core/internal/libs/commands"
	"servon/core/internal/libs/managers"
	"servon/core/internal/libs/utils"

	"github.com/spf13/cobra"
)

type CommandProvider struct {
	root *cobra.Command
	*utils.ShellUtil
	fullManager *managers.FullManager
	webServer   *utils.WebServer
}

func NewCommandProvider(fullManager *managers.FullManager, webServer *utils.WebServer) *CommandProvider {
	p := &CommandProvider{
		ShellUtil:   &utils.DefaultShellUtil,
		fullManager: fullManager,
		webServer:   webServer,
		root:        commands.RootCmd,
	}

	p.AddCommand(commands.GetDeployCommand())
	p.AddCommand(commands.GetServeCommand(p.webServer, p.fullManager))
	p.AddCommand(commands.GetVersionCommand(p.fullManager.VersionManager))
	p.AddCommand(commands.GetUpgradeCommand(p.fullManager.VersionManager))
	p.AddCommand(commands.GetSoftwareCommand(p.fullManager.SoftManager))
	p.AddCommand(commands.GetGitRootCommand(p.fullManager.GitManager))

	return p
}

// GetRootCommand 获取根命令
func (c *CommandProvider) GetRootCommand() *cobra.Command {
	return c.root
}

// AddCommand 添加命令
func (c *CommandProvider) AddCommand(command *cobra.Command) {
	c.root.AddCommand(command)
}
