package softwares

import (
	"fmt"
	"os/exec"
)

// Apt 提供apt包管理器的基本操作
type Apt struct {
	outputChan chan string
}

// NewApt 创建一个新的Apt实例
func NewApt(outputChan chan string) *Apt {
	return &Apt{
		outputChan: outputChan,
	}
}

// Update 更新软件包索引
func (a *Apt) Update() error {
	a.outputChan <- "更新软件包索引..."
	cmd := exec.Command("sudo", "apt-get", "update")
	output, err := cmd.CombinedOutput()
	if err != nil {
		a.outputChan <- fmt.Sprintf("更新索引失败:\n%s", string(output))
		return fmt.Errorf("更新索引失败: %v", err)
	}
	a.outputChan <- string(output)
	return nil
}

// Install 安装指定的软件包
func (a *Apt) Install(packages ...string) error {
	args := append([]string{"apt-get", "install", "-y"}, packages...)
	a.outputChan <- fmt.Sprintf("安装软件包: %v...", packages)
	cmd := exec.Command("sudo", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		a.outputChan <- fmt.Sprintf("安装失败:\n%s", string(output))
		return fmt.Errorf("安装失败: %v", err)
	}
	a.outputChan <- string(output)
	return nil
}

// Remove 移除指定的软件包
func (a *Apt) Remove(packages ...string) error {
	args := append([]string{"apt-get", "remove", "-y"}, packages...)
	a.outputChan <- fmt.Sprintf("移除软件包: %v...", packages)
	cmd := exec.Command("sudo", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		a.outputChan <- fmt.Sprintf("移除失败:\n%s", string(output))
		return fmt.Errorf("移除失败: %v", err)
	}
	a.outputChan <- string(output)
	return nil
}

// Purge 完全移除软件包及其配置文件
func (a *Apt) Purge(packages ...string) error {
	args := append([]string{"apt-get", "purge", "-y"}, packages...)
	a.outputChan <- "清理配置文件..."
	cmd := exec.Command("sudo", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		a.outputChan <- fmt.Sprintf("清理失败:\n%s", string(output))
		return fmt.Errorf("清理失败: %v", err)
	}
	a.outputChan <- string(output)
	return nil
}

// IsInstalled 检查软件包是否已安装
func (a *Apt) IsInstalled(packageName string) bool {
	cmd := exec.Command("dpkg", "-l", packageName)
	return cmd.Run() == nil
}
