package libs

import (
	"fmt"
	"os/exec"
)

type AptManager struct{}

func NewAptManager() *AptManager {
	return &AptManager{}
}

// AptUpdate 更新软件包索引
func (p *AptManager) AptUpdate() error {
	cmd := exec.Command("sudo", "apt-get", "update")
	if err := StreamCommand(cmd); err != nil {
		return fmt.Errorf("更新索引失败: %v", err)
	}
	return nil
}

// AptInstall 安装指定的软件包
func (p *AptManager) AptInstall(packages ...string) error {
	args := append([]string{"apt-get", "install", "-y"}, packages...)
	cmd := exec.Command("sudo", args...)

	if err := StreamCommand(cmd); err != nil {
		return fmt.Errorf("安装失败: %v", err)
	}

	InfoWithSpace("%s", fmt.Sprintf("安装成功: %v", packages))

	return nil
}

// AptRemove 移除指定的软件包
func (p *AptManager) AptRemove(packages ...string) error {
	args := append([]string{"apt-get", "remove", "-y"}, packages...)
	cmd := exec.Command("sudo", args...)
	if err := StreamCommand(cmd); err != nil {
		return fmt.Errorf("移除失败: %v", err)
	}
	return nil
}

// AptPurge 完全移除软件包及其配置文件
func (p *AptManager) AptPurge(packages ...string) error {
	args := append([]string{"apt-get", "purge", "-y"}, packages...)
	cmd := exec.Command("sudo", args...)
	if err := StreamCommand(cmd); err != nil {
		return fmt.Errorf("清理失败: %v", err)
	}
	return nil
}

// AptIsInstalled 检查软件包是否已安装
func (p *AptManager) AptIsInstalled(packageName string) bool {
	cmd := exec.Command("dpkg", "-l", packageName)
	return cmd.Run() == nil
}
