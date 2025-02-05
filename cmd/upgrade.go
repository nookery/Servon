package cmd

import (
	"servon/internal/upgrade"

	"github.com/spf13/cobra"
)

var (
	checkOnly bool
)

// UpgradeCmd 表示 upgrade 命令
var UpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "检查更新或升级到最新版本",
	Long: `检查是否有新版本可用，并可选择自动升级。

可选参数：
  --check-only: 仅检查是否有更新可用，不执行升级

示例：
  servon upgrade          # 检查并升级到最新版本
  servon upgrade --check  # 仅检查是否有更新可用`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return upgrade.CheckAndUpgrade(checkOnly)
	},
}

func init() {
	UpgradeCmd.Flags().BoolVar(&checkOnly, "check", false, "仅检查是否有更新可用")
}
