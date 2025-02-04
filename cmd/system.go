package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// SystemCmd 表示 system 命令
var SystemCmd = &cobra.Command{
	Use:   "system",
	Short: "系统管理",
	Long: `系统管理。

示例：
  servon system    # 系统管理`,
	Run: func(cmd *cobra.Command, args []string) {
		color.New(color.FgHiCyan).Printf("Servon system")
	},
}

func init() {
	// 版本命令不需要额外的参数配置
}
