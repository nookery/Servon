package system

import (
	"fmt"
	"os/exec"

	"servon/cmd/internal/utils"
)

// Apt 提供apt包管理器的基本操作
type Apt struct {
}

// NewApt 创建一个新的Apt实例
func NewApt() *Apt {
	return &Apt{}
}

// Update 更新软件包索引
func (a *Apt) Update() error {
	cmd := exec.Command("sudo", "apt-get", "update")
	if err := utils.StreamCommand(cmd); err != nil {
		return fmt.Errorf("更新索引失败: %v", err)
	}
	return nil
}

// Install 安装指定的软件包
func (a *Apt) Install(packages ...string) error {
	args := append([]string{"apt-get", "install", "-y"}, packages...)
	cmd := exec.Command("sudo", args...)

	if err := utils.StreamCommand(cmd); err != nil {
		return fmt.Errorf("安装失败: %v", err)
	}
	return nil
}

// Remove 移除指定的软件包
func (a *Apt) Remove(packages ...string) error {
	args := append([]string{"apt-get", "remove", "-y"}, packages...)
	cmd := exec.Command("sudo", args...)
	if err := utils.StreamCommand(cmd); err != nil {
		return fmt.Errorf("移除失败: %v", err)
	}
	return nil
}

// Purge 完全移除软件包及其配置文件
func (a *Apt) Purge(packages ...string) error {
	args := append([]string{"apt-get", "purge", "-y"}, packages...)
	cmd := exec.Command("sudo", args...)
	if err := utils.StreamCommand(cmd); err != nil {
		return fmt.Errorf("清理失败: %v", err)
	}
	return nil
}

// IsInstalled 检查软件包是否已安装
func (a *Apt) IsInstalled(packageName string) bool {
	cmd := exec.Command("dpkg", "-l", packageName)
	return cmd.Run() == nil
}
