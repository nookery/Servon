package libs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Version information
var (
	// 这些变量会在构建时通过 ldflags 注入
	Version    = "dev"     // 当前版本
	CommitHash = "none"    // Git commit hash
	BuildTime  = "unknown" // 构建时间
)

type VersionInfo struct {
	Version    string `json:"version"`
	CommitHash string `json:"commitHash"`
	BuildTime  string `json:"buildTime"`
}

type VersionManager struct {
	VersionInfo
	isDevVersion bool // 标记是否为开发版本
}

func NewVersionManager() *VersionManager {
	isDevVersion := false

	// 尝试从 package.json 读取版本(仅开发环境)
	if Version == "dev" {
		if ver := readPackageVersion(); ver != "" {
			Version = fmt.Sprintf("%s (dev)", ver)
			isDevVersion = true
		}
	}

	return &VersionManager{
		VersionInfo: VersionInfo{
			Version:    Version,
			CommitHash: CommitHash,
			BuildTime:  BuildTime,
		},
		isDevVersion: isDevVersion,
	}
}

// readPackageVersion 尝试从 package.json 读取版本号
func readPackageVersion() string {
	data, err := os.ReadFile("package.json")
	if err != nil {
		return ""
	}

	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return ""
	}

	return pkg.Version
}

// GetVersion 返回当前版本号
func (c *VersionManager) GetVersion() string {
	return c.Version
}

// GetVersionInfo 返回完整的版本信息
func (c *VersionManager) GetVersionInfo() VersionInfo {
	return c.VersionInfo
}

// GetLatestVersion 从 GitHub 获取最新发布版本
func (c *VersionManager) GetLatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/nookery/servon/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return strings.TrimPrefix(release.TagName, "v"), nil
}

// GetVersionCommand 返回版本命令
func (c *VersionManager) GetVersionCommand() *cobra.Command {
	return NewCommand(CommandOptions{
		Use:     "version",
		Short:   "显示版本信息",
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			DefaultPrinter.Printf("Version:     %s\n", c.Version)
			DefaultPrinter.Printf("Git Commit:  %s\n", c.CommitHash)
			DefaultPrinter.Printf("Built Time:  %s\n", c.BuildTime)
			if c.isDevVersion {
				DefaultPrinter.Printf("\n%s\n", "这是开发版本，版本号来自 package.json")
			}
		},
	})
}

// GetUpgradeCommand 返回升级命令
func (c *VersionManager) GetUpgradeCommand() *cobra.Command {
	return NewCommand(CommandOptions{
		Use:     "upgrade",
		Short:   "升级到最新版",
		Aliases: []string{"u", "up"},
		Run: func(cmd *cobra.Command, args []string) {
			DefaultPrinter.Printf("正在检查最新版本...\n")

			latestVersion, err := c.GetLatestVersion()
			if err != nil {
				DefaultPrinter.Printf("获取最新版本失败: %v\n", err)
				return
			}

			if latestVersion == c.Version {
				DefaultPrinter.Printf("当前已是最新版本: %s\n", c.Version)
				return
			}

			DefaultPrinter.Printf("发现新版本: %s，正在下载升级脚本...\n", latestVersion)

			resp, err := http.Get("https://raw.githubusercontent.com/nookery/servon/main/install.sh")
			if err != nil {
				DefaultPrinter.Printf("下载升级脚本失败: %v\n", err)
				return
			}
			defer resp.Body.Close()

			file, err := os.Create("install.sh")
			if err != nil {
				DefaultPrinter.Printf("创建升级脚本文件失败: %v\n", err)
				return
			}
			defer file.Close()

			_, err = io.Copy(file, resp.Body)
			if err != nil {
				DefaultPrinter.Printf("写入升级脚本文件失败: %v\n", err)
				return
			}

			DefaultPrinter.Printf("下载完成，正在执行升级脚本...\n")
			err = RunShell("bash", "install.sh")
			if err != nil {
				DefaultPrinter.Printf("执行升级脚本失败: %v\n", err)
				return
			}

			DefaultPrinter.Printf("升级完成！\n")
		},
	})
}
