package api

import (
	"fmt"
	"os/exec"
	"servon/core/libs"
	"servon/core/templates"

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

// CommandProvider 命令行命令执行器
type CommandProvider struct {
	RootCmd *cobra.Command
}

// AddCommand 添加命令
func (p *CommandProvider) AddCommand(cmd *cobra.Command) {
	p.RootCmd.AddCommand(cmd)
}

// NewCommand 创建一个标准化的命令
func (c *CommandApi) NewCommand(opts CommandOptions) *cobra.Command {
	setCustomErrPrefix := true
	setCustomUsageTemplate := true
	setCustomHelpFunc := true

	cmd := &cobra.Command{
		Use:           opts.Use,
		Short:         opts.Short,
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          opts.Args,
		PreRun: func(cmd *cobra.Command, args []string) {
			// libs.Infoln("🚀 开始执行命令 PreRun")
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// libs.Infoln("🚀 开始执行命令 PreRunE")
			return nil
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// libs.Infoln("🚀 开始执行命令 PersistentPreRun")
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// libs.Infoln("🚀 开始执行命令 PersistentPreRunE")
			return nil
		},
		Run: opts.Run,
		PostRun: func(cmd *cobra.Command, args []string) {
			// libs.Infoln("🎉 命令执行成功 PostRun")
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			// libs.Infoln("🎉 命令执行完成 PostRunE")
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			// libs.Infoln("🎉 命令执行完成 PersistentPostRun")
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			// libs.Infoln("🎉 命令执行完成 PersistentPostRunE")
			return nil
		},
	}

	if setCustomErrPrefix {
		cmd.SetErrPrefix("❌ 发生了错误")
	}

	// 自定义错误处理
	cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		c.Printf("%s\n", red("❌ 错误："+fmt.Sprintf("%v", err)))
		return nil
	})

	// 自定义帮助
	if setCustomHelpFunc {
		cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
			c.Printf("%s\n", purple("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
			c.Printf("📌 命令: %s\n", yellow(c.Use))
			c.Printf("📝 描述: %s\n", green(c.Short))
			c.Printf("\n%s\n", blue("🎯 参数列表:"))
			c.Printf("%s\n", blue(c.LocalFlags().FlagUsages()))
			c.Printf("%s\n", cyan("🎯 可用命令:"))
			for _, command := range c.Commands() {
				c.Printf("  %-35s %s\n", cyan(command.Use), command.Short)
			}
			c.Printf("%s\n", purple("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
		})
	}

	// 自定义使用说明模板
	if setCustomUsageTemplate {
		cmd.SetUsageTemplate(templates.UsageTemplate())
	}

	return cmd
}

func (c *CommandApi) AddCommand(cmd *cobra.Command) {
	c.rootCmd.AddCommand(cmd)
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
