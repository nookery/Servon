package cmd

import (
	"fmt"
	"servon/internal/version"

	"github.com/fatih/color"
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
		// 获取当前版本
		currentVersion := version.GetVersion()

		// 检查最新版本
		latestVersion, err := version.GetLatestVersion()
		if err != nil {
			return fmt.Errorf("检查更新失败: %v", err)
		}

		// 比较版本
		needsUpgrade, err := version.NeedsUpgrade(currentVersion, latestVersion)
		if err != nil {
			return fmt.Errorf("版本比较失败: %v", err)
		}

		if !needsUpgrade {
			color.Green("当前已是最新版本 %s", currentVersion)
			return nil
		}

		color.Yellow("发现新版本: %s (当前版本: %s)", latestVersion, currentVersion)

		if checkOnly {
			color.Cyan("运行 'servon upgrade' 命令来升级到最新版本")
			return nil
		}

		// 执行升级
		color.Cyan("正在升级到版本 %s ...", latestVersion)
		if err := version.DoUpgrade(); err != nil {
			return fmt.Errorf("升级失败: %v", err)
		}

		color.Green("升级成功！")
		return nil
	},
}

func init() {
	UpgradeCmd.Flags().BoolVar(&checkOnly, "check", false, "仅检查是否有更新可用")
}
