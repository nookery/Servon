package libs

import (
	"fmt"
	"os/exec"
	"strings"
)

// Dpkg 提供dpkg包管理器的基本操作
type Dpkg struct {
	outputChan chan string
}

// NewDpkg 创建一个新的Dpkg实例
func NewDpkg(outputChan chan string) *Dpkg {
	return &Dpkg{
		outputChan: outputChan,
	}
}

// IsInstalled 检查软件包是否已安装
// 返回 true 表示已安装，false 表示未安装
func (d *Dpkg) IsInstalled(packageName string) bool {
	cmd := exec.Command("dpkg", "-l", packageName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	// dpkg -l 的输出中，如果包已安装会显示 "ii  package"
	return strings.Contains(string(output), fmt.Sprintf("ii  %s", packageName))
}

// GetVersion 获取已安装软件包的版本
// 如果软件包未安装，返回空字符串
func (d *Dpkg) GetVersion(packageName string) string {
	cmd := exec.Command("dpkg", "-l", packageName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	// 解析输出获取版本信息
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ii  "+packageName) {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				return fields[2]
			}
		}
	}
	return ""
}

// ListPackages 列出所有已安装的软件包
func (d *Dpkg) ListPackages() ([]string, error) {
	cmd := exec.Command("dpkg", "--get-selections")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("获取软件包列表失败: %v", err)
	}

	var packages []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "install") {
			pkg := strings.Fields(line)[0]
			packages = append(packages, pkg)
		}
	}
	return packages, nil
}

// GetArchitecture 获取系统架构
func (d *Dpkg) GetArchitecture() (string, error) {
	cmd := exec.Command("dpkg", "--print-architecture")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取系统架构失败: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}
