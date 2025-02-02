package softwares

import (
	"fmt"
	"os/exec"
	"strings"
)

type Nginx struct {
	info SoftwareInfo
}

func NewNginx() *Nginx {
	return &Nginx{
		info: SoftwareInfo{
			Name:        "nginx",
			Description: "高性能的 HTTP 和反向代理服务器",
		},
	}
}

func (n *Nginx) Install() (chan string, error) {
	outputChan := make(chan string, 100)

	go func() {
		defer close(outputChan)

		// 更新软件包索引
		outputChan <- "更新软件包索引..."
		updateCmd := exec.Command("sudo", "apt-get", "update")
		if output, err := updateCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("更新索引失败:\n%s", string(output))
			return
		}

		// 安装 Nginx
		outputChan <- "安装 Nginx..."
		installCmd := exec.Command("sudo", "apt-get", "install", "-y", "nginx")
		if output, err := installCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("安装失败:\n%s", string(output))
			return
		}

		// 启动服务
		outputChan <- "启动 Nginx 服务..."
		startCmd := exec.Command("sudo", "systemctl", "start", "nginx")
		if output, err := startCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("启动服务失败:\n%s", string(output))
			return
		}

		outputChan <- "Nginx 安装完成"
	}()

	return outputChan, nil
}

func (n *Nginx) Uninstall() (chan string, error) {
	outputChan := make(chan string, 100)

	go func() {
		defer close(outputChan)

		// 停止服务
		outputChan <- "停止 Nginx 服务..."
		stopCmd := exec.Command("sudo", "systemctl", "stop", "nginx")
		if output, err := stopCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("停止服务失败:\n%s", string(output))
		}

		// 卸载软件
		outputChan <- "卸载 Nginx..."
		removeCmd := exec.Command("sudo", "apt-get", "remove", "-y", "nginx", "nginx-common")
		if output, err := removeCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("卸载失败:\n%s", string(output))
			return
		}

		// 清理配置文件
		outputChan <- "清理配置文件..."
		purgeCmd := exec.Command("sudo", "apt-get", "purge", "-y", "nginx", "nginx-common")
		if output, err := purgeCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("清理失败:\n%s", string(output))
			return
		}

		outputChan <- "Nginx 卸载完成"
	}()

	return outputChan, nil
}

func (n *Nginx) GetStatus() (map[string]string, error) {
	// 检查是否安装
	cmd := exec.Command("dpkg", "-l", "nginx")
	if err := cmd.Run(); err != nil {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	// 检查服务状态
	statusCmd := exec.Command("systemctl", "status", "nginx")
	output, err := statusCmd.CombinedOutput()
	status := "stopped"
	if err == nil && strings.Contains(string(output), "Active: active (running)") {
		status = "running"
	}

	// 获取版本
	version := ""
	verCmd := exec.Command("nginx", "-v")
	if verOutput, err := verCmd.CombinedOutput(); err == nil {
		version = strings.TrimSpace(string(verOutput))
	}

	return map[string]string{
		"status":  status,
		"version": version,
	}, nil
}

func (n *Nginx) Stop() error {
	cmd := exec.Command("sudo", "systemctl", "stop", "nginx")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("停止服务失败: %v\n%s", err, string(output))
	}
	return nil
}

func (n *Nginx) GetInfo() SoftwareInfo {
	return n.info
}
