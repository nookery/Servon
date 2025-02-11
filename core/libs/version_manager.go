package libs

import (
	"encoding/json"
	"fmt"
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
			printer.Printf("Version:     %s\n", c.Version)
			printer.Printf("Git Commit:  %s\n", c.CommitHash)
			printer.Printf("Built Time:  %s\n", c.BuildTime)
			if c.isDevVersion {
				printer.Printf("\n%s\n", "这是开发版本，版本号来自 package.json")
			}
		},
	})
}
