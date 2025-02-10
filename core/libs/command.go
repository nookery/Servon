package libs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

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

// PrintCommandHelp 打印标准格式的命令帮助信息
func PrintCommandHelp(cmd *cobra.Command) {
	fmt.Println()

	// 首先显示 Long 描述（包含 ASCII 艺术和描述文本）
	if cmd.Long != "" {
		// 如果是多行，原样输出
		if strings.Contains(cmd.Long, "\n") {
			fmt.Println(cmd.Long)
		} else {
			fmt.Println(color.New(color.BgGreen).Sprintf(" ✨ %s ✨ ", cmd.Long))
		}
	}

	// 自动获取所有子命令及其描述
	commands := make(map[string]string)
	for _, subCmd := range cmd.Commands() {
		if !subCmd.Hidden {
			commands[subCmd.Name()] = subCmd.Short
		}
	}

	// 使用方法
	color.New(color.FgHiWhite).Printf("\n📌 使用方法: ")
	color.New(color.FgCyan).Printf("%s\n\n", cmd.UseLine())

	// 添加参数列表展示
	if cmd.HasFlags() {
		color.New(color.FgHiWhite).Println("🎯 参数选项:")
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			// 构建默认值字符串
			defaultValue := ""
			if flag.DefValue != "" {
				defaultValue = fmt.Sprintf("(默认值: %s)", flag.DefValue)
			}

			// 构建参数名称
			name := ""
			if flag.Shorthand != "" && flag.Shorthand != flag.Name {
				name = fmt.Sprintf("-%s, --%s", flag.Shorthand, flag.Name)
			} else {
				name = fmt.Sprintf("--%s", flag.Name)
			}

			color.New(color.FgCyan).Printf("  ▶️  %-20s", name)
			color.New(color.FgWhite).Printf("%s %s\n", flag.Usage, defaultValue)
		})
	}

	// 可用命令列表
	if len(commands) > 0 {
		color.New(color.FgHiWhite).Printf("\n🔧 可用命令:\n")
		for name, desc := range commands {
			color.New(color.FgCyan).Printf("  ▶️  %s", name)
			color.New(color.FgWhite).Printf("\t%s\n", desc)
		}
	}

	fmt.Println()
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
