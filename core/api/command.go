package api

import (
	"os/exec"
	"servon/core/libs"

	"github.com/spf13/cobra"
)

type Command struct {
	commandProvider libs.CommandProvider
}

func NewCommandApi() Command {
	return Command{
		commandProvider: libs.NewCommandProvider(),
	}
}

func (c *Command) AddCommand(cmd *cobra.Command) {
	c.commandProvider.AddCommand(cmd)
}

// CheckCommandArgs 检查命令参数
func (c *Command) CheckCommandArgs(cmd *cobra.Command, args []string) error {
	return libs.CheckCommandArgs(cmd, args)
}

func (c *Command) GetRootCommand() *cobra.Command {
	return c.commandProvider.RootCmd
}

// StreamCommand 执行命令并打印输出
func (c *Command) StreamCommand(cmd *exec.Cmd) error {
	return libs.StreamCommand(cmd)
}

func (c *Command) RunShell(command string, args ...string) error {
	return libs.Execute(command, args...)
}

func (c *Command) RunShellWithOutput(command string, args ...string) (string, error) {
	return libs.ExecuteWithOutput(command, args...)
}

// PrintCommandHelp 打印命令帮助
func (c *Command) PrintCommandHelp(cmd *cobra.Command) {
	libs.PrintCommandHelp(cmd)
}
