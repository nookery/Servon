package cmd

import (
	"servon/internal/version"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// VersionCmd 表示 version 命令
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long: `显示 Servon 的版本信息。

示例：
  servon version    # 显示当前版本`,
	Run: func(cmd *cobra.Command, args []string) {
		color.New(color.FgHiCyan).Printf("Servon version %s\n", version.GetVersion())
	},
}

func init() {
	// 版本命令不需要额外的参数配置
}
