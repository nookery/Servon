package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// SoftwareCmd 表示 software 命令
var SoftwareCmd = &cobra.Command{
	Use:   "software",
	Short: "软件管理",
	Long: `软件管理。

示例：
  servon software    # 软件管理`,
	Run: func(cmd *cobra.Command, args []string) {
		color.New(color.FgHiCyan).Printf("Servon software")
	},
}

func init() {
	// 版本命令不需要额外的参数配置
}
