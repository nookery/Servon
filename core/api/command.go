package api

import (
	"os/exec"
	"servon/core/libs"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Command struct {
	rootCmd *cobra.Command
}

type CommandOptions = libs.CommandOptions

var (
	titleColor = color.New(color.FgHiCyan, color.Bold)
	infoColor  = color.New(color.FgHiWhite)
)

// CommandProvider 命令行命令执行器
type CommandProvider struct {
	RootCmd *cobra.Command
}

// AddCommand 添加命令
func (p *CommandProvider) AddCommand(cmd *cobra.Command) {
	p.RootCmd.AddCommand(cmd)
}

func NewCommandApi() Command {
	rootCmd := libs.NewCommand(CommandOptions{
		Use:   "servon",
		Short: "Servon - A lightweight server management panel",
	})

	return Command{
		rootCmd: rootCmd,
	}
}

func (c *Command) AddCommand(cmd *cobra.Command) {
	c.rootCmd.AddCommand(cmd)
}

// CheckCommandArgs 检查命令参数
func (c *Command) CheckCommandArgs(cmd *cobra.Command, args []string) error {
	return libs.CheckCommandArgs(cmd, args)
}

func (c *Command) GetRootCommand() *cobra.Command {
	return c.rootCmd
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

// NewCommand 创建一个标准化的命令
func (c *Command) NewCommand(opts CommandOptions) *cobra.Command {
	return libs.NewCommand(opts)
}
