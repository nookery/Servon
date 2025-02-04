package utils

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// PrintCommandHelp 打印标准格式的命令帮助信息
func PrintCommandHelp(cmd *cobra.Command, commands map[string]string) {
	// 使用方法
	color.New(color.FgHiWhite).Printf("使用方法:\t")
	color.New(color.FgCyan).Printf("servon %s COMMAND\n", cmd.Name())

	// 命令说明
	color.New(color.FgHiCyan).Printf("\nServon %s管理命令\n\n", cmd.Name())

	// 可用命令列表
	color.New(color.FgHiWhite).Println("可用命令:")
	for name, desc := range commands {
		color.New(color.FgCyan).Printf("  %s", name)
		color.New(color.FgWhite).Printf("\t%s\n", desc)
	}

	// 帮助提示
	color.New(color.FgHiWhite).Printf("\n使用 ")
	color.New(color.FgYellow).Printf("\"servon %s COMMAND --help\"", cmd.Name())
	color.New(color.FgHiWhite).Println(" 了解更多命令信息")
}
