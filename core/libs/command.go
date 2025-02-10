package libs

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type CommandOptions struct {
	Use   string
	Short string
	Run   func(cmd *cobra.Command, args []string)
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

// PrintList 打印列表
func PrintList(list []string, title string) {
	fmt.Println()
	color.New(color.FgHiCyan).Println(title)
	if len(list) == 0 {
		color.New(color.FgYellow).Println("  暂无数据")
		fmt.Println()
		return
	}
	for _, item := range list {
		color.New(color.FgCyan).Printf("  ▶️  %s\n", item)
	}
	fmt.Println()
}

// PrintError 打印错误信息
func PrintError(err error) {
	fmt.Println()
	color.New(color.FgHiRed).Printf("❌ 错误: %s\n", err.Error())
	fmt.Println()
}

// NewCommand 创建一个标准化的命令
func NewCommand(opts CommandOptions) *cobra.Command {
	setCustomErrPrefix := false
	setCustomUsageTemplate := false
	setCustomHelpFunc := false

	cmd := &cobra.Command{
		Use:           opts.Use,
		Short:         opts.Short,
		SilenceErrors: false,
		SilenceUsage:  false,
		PreRun: func(cmd *cobra.Command, args []string) {
			Infoln("🚀 开始执行命令 PreRun")
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			Infoln("🚀 开始执行命令 PreRunE")
			return nil
		},
		Run: opts.Run,
		PostRun: func(cmd *cobra.Command, args []string) {
			Infoln("🎉 命令执行成功 PostRun")
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			Infoln("🎉 命令执行完成 PostRunE")
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			Infoln("🎉 命令执行完成 PersistentPostRun")
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			Infoln("🎉 命令执行完成 PersistentPostRunE")
			return nil
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			Infoln("🚀 开始执行命令 PersistentPreRun")
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			Infoln("🚀 开始执行命令 PersistentPreRunE")
			return nil
		},
	}

	if setCustomErrPrefix {
		cmd.SetErrPrefix("❌ 发生了错误")
	}

	// 自定义错误处理
	cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		c.Printf("\x1b[1;31m❌ 错误：缺少必需的参数\x1b[0m\n")
		return nil
	})

	// 自定义帮助
	if setCustomHelpFunc {
		cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
			c.Printf("\x1b[1;36m🌈 命令帮助\x1b[0m\n")
			c.Printf("\x1b[1;35m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\x1b[0m\n")
			c.Printf("\x1b[1;33m📌 命令: %s\x1b[0m\n", c.Use)
			c.Printf("\x1b[1;32m📝 描述: %s\x1b[0m\n", c.Short)
			c.Printf("\x1b[1;34m🎯 参数列表:\x1b[0m\n")
			c.Printf("\x1b[1;34m%s\x1b[0m\n", c.LocalFlags().FlagUsages())
			c.Printf("\x1b[1;34m🎯 可用命令:\x1b[0m\n")
			for _, command := range c.Commands() {
				c.Printf("\x1b[1;34m  %s\x1b[0m %s\n", command.Use, command.Short)
			}
			c.Printf("\n\x1b[1;36m✨ 示例:\x1b[0m\n")
			c.Printf("\x1b[1;36m%s [参数]\x1b[0m\n", c.CommandPath())
		})
	}

	// 自定义使用说明模板
	if setCustomUsageTemplate {
		cmd.SetUsageTemplate(`
	` + "\x1b[1;36m" + `🌈 命令说明` + "\x1b[0m" + `
	` + "\x1b[1;35m" + `━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━` + "\x1b[0m" + `
	` + "\x1b[1;33m" + `📌 命令:` + "\x1b[0m" + ` {{.UseLine}}
	` + "\x1b[1;32m" + `📝 描述:` + "\x1b[0m" + ` {{.Short}}

	` + "\x1b[1;34m" + `🎯 参数列表:` + "\x1b[0m" + `
	{{.LocalFlags.FlagUsages}}
	` + "\x1b[1;36m" + `✨ 示例:` + "\x1b[0m" + `{{.CommandPath}} [参数]

	` + "\x1b[1;33m" + `💡 提示:` + "\x1b[0m" + ` 使用 -h 或 --help 查看更多帮助信息
	` + "\x1b[1;35m" + `━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━` + "\x1b[0m\n\n" + `
	`)
	}

	// 确保错误不会传播到父命令
	if cmd.Root() != nil {
		cmd.Root().SilenceErrors = true
	}

	return cmd
}
