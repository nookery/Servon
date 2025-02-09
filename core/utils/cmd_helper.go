package utils

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// PrintCommandHelp 打印标准格式的命令帮助信息
func PrintCommandHelp(cmd *cobra.Command) {
	fmt.Println()

	// 首先显示 Long 描述（包含 ASCII 艺术和描述文本）
	fmt.Println(cmd.Long)

	// 自动获取所有子命令及其描述
	commands := make(map[string]string)
	for _, subCmd := range cmd.Commands() {
		if !subCmd.Hidden {
			commands[subCmd.Name()] = subCmd.Short
		}
	}

	// 使用方法
	color.New(color.FgHiWhite).Printf("\n📌 使用方法: ")
	color.New(color.FgCyan).Printf("%s\n", cmd.UseLine())

	// 可用命令列表
	color.New(color.FgHiWhite).Println("🔧 可用命令:")
	for name, desc := range commands {
		color.New(color.FgCyan).Printf("  ▶️  %s", name)
		color.New(color.FgWhite).Printf("\t%s\n", desc)
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
