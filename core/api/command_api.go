package api

import (
	"os/exec"
	"servon/core/libs"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type CommandApi struct {
	rootCmd *cobra.Command
}

type CommandOptions = libs.CommandOptions

// 定义颜色打印函数
var (
	cyan   = color.New(color.FgCyan).SprintFunc()
	purple = color.New(color.FgMagenta).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
)

func NewCommandApi() CommandApi {
	api := CommandApi{}

	api.rootCmd = api.NewCommand(CommandOptions{
		Use:   "servon",
		Short: "Servon 是一个用于管理服务器的命令行工具",
	})

	return api
}

// AddCommand 添加命令
func (p *CommandApi) AddCommand(cmd *cobra.Command) {
	p.rootCmd.AddCommand(cmd)
}

// CheckCommandArgs 检查命令参数
func (c *CommandApi) CheckCommandArgs(cmd *cobra.Command, args []string) error {
	return libs.CheckCommandArgs(cmd, args)
}

func (c *CommandApi) GetRootCommand() *cobra.Command {
	return c.rootCmd
}

// StreamCommand 执行命令并打印输出
func (c *CommandApi) StreamCommand(cmd *exec.Cmd) error {
	return libs.StreamCommand(cmd)
}

func (c *CommandApi) RunShell(command string, args ...string) error {
	return libs.Execute(command, args...)
}

func (c *CommandApi) RunShellWithOutput(command string, args ...string) (string, error) {
	return libs.ExecuteWithOutput(command, args...)
}

func (c *CommandApi) NewCommand(opts CommandOptions) *cobra.Command {
	return libs.NewCommand(opts)
}
