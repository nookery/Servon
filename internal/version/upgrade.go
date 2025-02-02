package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

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
