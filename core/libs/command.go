package libs

import (
	"fmt"
	"os"
	"os/exec"
	"servon/core/templates"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type CommandOptions struct {
	Use     string
	Short   string
	Args    cobra.PositionalArgs
	Run     func(cmd *cobra.Command, args []string)
	Aliases []string
}

// CheckCommandArgs 检查命令参数
func CheckCommandArgs(cmd *cobra.Command, args []string) error {
	// 如果命令没有设置 Args 要求，则至少需要一个参数
	if cmd.Args == nil {
		if len(args) == 0 {
			return fmt.Errorf("至少需要一个参数")
		}
		return nil
	}

	// 使用命令自带的参数验证
	err := cmd.Args(cmd, args)
	if err != nil {
		return fmt.Errorf("参数验证失败: %v", err)
	}

	Info("参数验证成功")
	return nil
}

func Execute(command string, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("command is required")
	}

	// 使用青色（Cyan）输出命令和参数，用空格连接参数
	color.Cyan("📺 %s %s", command, joinArgs(args))

	execCmd := exec.Command(command, args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	return execCmd.Run()
}

func ExecuteWithOutput(command string, args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("command is required")
	}

	execCmd := exec.Command(command, args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	output, err := execCmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// joinArgs 将参数数组连接成字符串，去掉方括号
func joinArgs(args []string) string {
	result := ""
	for i, arg := range args {
		if i > 0 {
			result += " "
		}
		result += arg
	}
	return result
}

// NewCommand 创建一个标准化的命令
func NewCommand(opts CommandOptions) *cobra.Command {
	setCustomErrPrefix := true
	setCustomUsageTemplate := true
	setCustomHelpFunc := true

	cmd := &cobra.Command{
		Use:           opts.Use,
		Short:         opts.Short,
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          opts.Args,
		Aliases:       opts.Aliases,
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
		c.Printf("%s\n", color.New(color.FgHiRed).Sprintf("%s", "❌ 错误："+fmt.Sprintf("%v", err)))
		return nil
	})

	// 自定义帮助
	if setCustomHelpFunc {
		cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
			c.Printf("%s\n", color.New(color.FgHiCyan).Sprintf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
			c.Printf("📌 命令: %s\n", color.New(color.FgHiYellow).Sprintf(c.Use))
			c.Printf("📝 描述: %s\n", color.New(color.FgHiGreen).Sprintf(c.Short))
			c.Printf("\n%s\n", color.New(color.FgHiBlue).Sprintf("🎯 参数列表:"))
			c.Printf("%s\n", color.New(color.FgHiCyan).Sprintf("%s", c.LocalFlags().FlagUsages()))
			c.Printf("%s\n", color.New(color.FgHiCyan).Sprintf("🎯 可用命令:"))
			for _, command := range c.Commands() {
				alias := ""
				if len(command.Aliases) > 0 {
					alias = "(" + joinArgs(command.Aliases) + ")"
				}

				nameAndAlias := ""
				if alias != "" {
					nameAndAlias = fmt.Sprintf("%s %s", command.Use, alias)
				} else {
					nameAndAlias = command.Use
				}
				c.Printf("  %-35s%s\n", color.New(color.FgHiCyan).Sprintf("%s", nameAndAlias), command.Short)
			}
			c.Printf("%s\n", color.New(color.FgHiCyan).Sprintf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
		})
	}

	// 自定义使用说明模板
	if setCustomUsageTemplate {
		cmd.SetUsageTemplate(templates.UsageTemplate())
	}

	return cmd
}
