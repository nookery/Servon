package libs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

// Version information
var (
	Version   = "0.1.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// VersionProvider 版本信息提供者
type VersionProvider struct {
	version   string
	buildTime string
	gitCommit string
}

// NewVersionProvider 创建版本信息提供者
func NewVersionProvider() VersionProvider {
	return VersionProvider{
		version:   Version,
		buildTime: BuildTime,
		gitCommit: GitCommit,
	}
}

// GetVersion 获取版本信息
func (p *VersionProvider) GetVersion() string {
	return p.version
}

// GetFullVersionInfo 获取完整版本信息
func (p *VersionProvider) GetFullVersionInfo() string {
	return fmt.Sprintf("Version: %s\nBuild Time: %s\nGit Commit: %s",
		p.version,
		p.buildTime,
		p.gitCommit,
	)
}

// CheckAndUpgrade 检查并执行升级
func (v *VersionProvider) CheckAndUpgrade(checkOnly bool) error {
	// 获取当前版本
	currentVersion := v.GetVersion()

	// 检查最新版本
	latestVersion, err := v.GetLatestVersion()
	if err != nil {
		return fmt.Errorf("检查更新失败: %v", err)
	}

	// 比较版本
	needsUpgrade, err := v.NeedsUpgrade(currentVersion, latestVersion)
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
	if err := v.DoUpgrade(); err != nil {
		return fmt.Errorf("升级失败: %v", err)
	}

	color.Green("升级成功！")
	return nil
}

// GetLatestVersion 获取最新版本
func (v *VersionProvider) GetLatestVersion() (string, error) {
	// TODO: 实现获取最新版本的逻辑
	return "", nil
}

// NeedsUpgrade 检查是否需要升级
func (v *VersionProvider) NeedsUpgrade(currentVersion, latestVersion string) (bool, error) {
	// TODO: 实现版本比较逻辑
	return false, nil
}

// DoUpgrade 执行升级
func (v *VersionProvider) DoUpgrade() error {
	// TODO: 实现升级逻辑
	return nil
}

// GetLatestVersion 从 GitHub 获取最新版本
func GetLatestVersion() (string, error) {
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

// NeedsUpgrade 检查是否需要升级
func NeedsUpgrade(current, latest string) (bool, error) {
	// 如果是开发版本，总是建议升级
	if strings.Contains(current, "development") {
		return true, nil
	}

	// 移除版本号前的 'v' 前缀
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	// 简单的版本比较
	return latest > current, nil
}

// DoUpgrade 执行升级
func DoUpgrade() error {
	// 获取安装脚本
	resp, err := http.Get("https://raw.githubusercontent.com/nookery/servon/main/install.sh")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 保存安装脚本
	tmpFile := fmt.Sprintf("/tmp/servon-install-%d.sh", os.Getpid())
	f, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	if _, err := f.ReadFrom(resp.Body); err != nil {
		f.Close()
		return err
	}
	f.Close()

	// 执行安装脚本
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		return fmt.Errorf("Windows 系统暂不支持自动升级，请访问 https://github.com/nookery/servon/releases 手动升级")
	} else {
		cmd = exec.Command("bash", tmpFile)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
