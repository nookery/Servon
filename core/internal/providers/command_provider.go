package providers

import (
	"servon/core/internal/commands"
	"servon/core/internal/managers"
	"servon/components/shell_util"
	"servon/components/web_server_util"

	"github.com/spf13/cobra"
)

type CommandProvider struct {
	root *cobra.Command
	*shell_util.ShellUtil
	fullManager *managers.FullManager
	webServer   *web_server_util.WebServer
}

func NewCommandProvider(fullManager *managers.FullManager, webServer *web_server_util.WebServer) *CommandProvider {
	p := &CommandProvider{
		ShellUtil:   &shell_util.DefaultShellUtil,
		fullManager: fullManager,
		webServer:   webServer,
		root:        commands.RootCmd,
	}

	p.AddCommand(commands.GetDeployCommand(p.fullManager))
	p.AddCommand(commands.GetServerCommand(p.webServer, p.fullManager))
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
