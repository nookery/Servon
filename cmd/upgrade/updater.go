package upgrade

import (
	"fmt"
	"servon/cmd/version"

	"github.com/fatih/color"
)

func CheckAndUpgrade(checkOnly bool) error {
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
}
