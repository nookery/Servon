package managers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"servon/components/utils"
	"strings"
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
	IsDevVersion bool // 标记是否为开发版本
}

func NewVersionManager() *VersionManager {
	isDevVersion := false

	// 尝试从 package.json 读取版本(仅开发环境)
	if Version == "dev" {
		if ver := utils.ReadPackageVersion(); ver != "" {
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
		IsDevVersion: isDevVersion,
	}
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
