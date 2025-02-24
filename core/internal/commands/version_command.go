package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"servon/core/internal/managers"

	"github.com/spf13/cobra"
)

// GetVersionCommand 返回版本命令
func GetVersionCommand(c *managers.VersionManager) *cobra.Command {
	return NewCommand(CommandOptions{
		Use:     "version",
		Short:   "显示版本信息",
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			PrintKeyValues(map[string]string{
				"Version":      c.Version,
				"CommitHash":   c.CommitHash,
				"BuildTime":    c.BuildTime,
				"IsDevVersion": fmt.Sprintf("%t", c.IsDevVersion),
			})
		},
	})
}

// GetUpgradeCommand 返回升级命令
func GetUpgradeCommand(c *managers.VersionManager) *cobra.Command {
	return NewCommand(CommandOptions{
		Use:     "upgrade",
		Short:   "升级到最新版",
		Aliases: []string{"u", "up"},
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("正在检查最新版本...")

			latestVersion, err := c.GetLatestVersion()
			if err != nil {
				logger.Errorf("获取最新版本失败: %v", err)
				return
			}

			if latestVersion == c.Version {
				logger.Infof("当前已是最新版本: %s", c.Version)
				return
			}

			logger.Infof("发现新版本: %s，正在下载升级脚本...", latestVersion)

			resp, err := http.Get("https://raw.githubusercontent.com/nookery/servon/main/install.sh")
			if err != nil {
				logger.Errorf("下载升级脚本失败: %v", err)
				return
			}
			defer resp.Body.Close()

			file, err := os.Create("install.sh")
			if err != nil {
				logger.Errorf("创建升级脚本文件失败: %v", err)
				return
			}
			defer file.Close()

			_, err = io.Copy(file, resp.Body)
			if err != nil {
				logger.Errorf("写入升级脚本文件失败: %v", err)
				return
			}

			logger.Info("下载完成，正在执行升级脚本...")
			err, _ = RunShell("bash", "install.sh")
			if err != nil {
				logger.Errorf("执行升级脚本失败: %v", err)
				return
			}

			logger.Info("升级完成！")
		},
	})
}
