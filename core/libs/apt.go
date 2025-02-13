package libs

import (
	"fmt"
	"strings"
)

type AptManager struct{}

func NewAptManager() *AptManager {
	return &AptManager{}
}

// AptUpdate 更新软件包索引
func (p *AptManager) AptUpdate() error {
	PrintInfo("正在更新软件包索引...")
	if err := RunShell("sudo", "apt-get", "update"); err != nil {
		return fmt.Errorf("更新索引失败: %v", err)
	}
	PrintSuccess("软件包索引更新成功")
	return nil
}

// AptInstall 安装指定的软件包
func (p *AptManager) AptInstall(packages ...string) error {
	if err := RunShell("sudo", "apt-get", "install", "-y", strings.Join(packages, " ")); err != nil {
		return fmt.Errorf("安装失败: %v", err)
	}

	DefaultPrinter.PrintInfo(fmt.Sprintf("Apt 安装成功: %v", packages))

	return nil
}

// AptRemove 移除指定的软件包
func (p *AptManager) AptRemove(packages ...string) error {
	if err := RunShell("sudo", "apt-get", "remove", "-y", strings.Join(packages, " ")); err != nil {
		return fmt.Errorf("移除失败: %v", err)
	}
	return nil
}

// AptPurge 完全移除软件包及其配置文件
func (p *AptManager) AptPurge(packages ...string) error {
	if err := RunShell("sudo", "apt-get", "purge", "-y", strings.Join(packages, " ")); err != nil {
		return fmt.Errorf("清理失败: %v", err)
	}
	return nil
}

// AptIsInstalled 检查软件包是否已安装
func (p *AptManager) AptIsInstalled(packageName string) bool {
	return RunShell("dpkg", "-l", packageName) == nil
}
