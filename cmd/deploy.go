package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// DeployCmd 表示 deploy 命令
var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "部署项目",
	Long: `部署项目。

示例：
  servon deploy    # 部署项目`,
	Run: func(cmd *cobra.Command, args []string) {
		color.New(color.FgHiCyan).Printf("Servon deploy")
	},
}

func init() {
	// 版本命令不需要额外的参数配置
}
