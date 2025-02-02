package softwares

import (
	"fmt"
	"os/exec"
	"strings"
)

type MySQL struct {
	info SoftwareInfo
}

func NewMySQL() *MySQL {
	return &MySQL{
		info: SoftwareInfo{
			Name:        "mysql",
			Description: "流行的关系型数据库",
		},
	}
}

func (m *MySQL) Install() (chan string, error) {
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

		// 安装 MySQL
		outputChan <- "安装 MySQL..."
		installCmd := exec.Command("sudo", "apt-get", "install", "-y", "mysql-server")
		if output, err := installCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("安装失败:\n%s", string(output))
			return
		}

		// 设置 root 密码为空（开发环境）
		outputChan <- "配置 MySQL..."
		configCmd := exec.Command("sudo", "mysql", "-e", "ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY '';")
		if output, err := configCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("配置失败:\n%s", string(output))
			return
		}

		outputChan <- "MySQL 安装完成"
	}()

	return outputChan, nil
}

func (m *MySQL) Uninstall() (chan string, error) {
	outputChan := make(chan string, 100)

	go func() {
		defer close(outputChan)

		// 停止服务
		outputChan <- "停止 MySQL 服务..."
		stopCmd := exec.Command("sudo", "systemctl", "stop", "mysql")
		output, err := stopCmd.CombinedOutput()
		if err != nil {
			outputChan <- fmt.Sprintf("停止服务失败:\n%s", string(output))
		}
		outputChan <- string(output)

		// 卸载软件
		outputChan <- "卸载 MySQL..."
		removeCmd := exec.Command("sudo", "apt-get", "remove", "-y", "mysql-server*")
		if output, err := removeCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("卸载失败:\n%s", string(output))
			return
		}
		outputChan <- string(output)

		// 清理配置文件
		outputChan <- "清理配置文件..."
		purgeCmd := exec.Command("sudo", "apt-get", "purge", "-y", "mysql-server*")
		if output, err := purgeCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("清理失败:\n%s", string(output))
			return
		}
		outputChan <- string(output)

		outputChan <- "MySQL 卸载完成"
	}()

	return outputChan, nil
}

func (m *MySQL) GetStatus() (map[string]string, error) {
	// 检查是否安装
	cmd := exec.Command("dpkg", "-l", "mysql-server")
	if err := cmd.Run(); err != nil {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	// 检查服务状态
	statusCmd := exec.Command("systemctl", "status", "mysql")
	output, err := statusCmd.CombinedOutput()
	status := "stopped"
	if err == nil && strings.Contains(string(output), "Active: active (running)") {
		status = "running"
	}

	// 获取版本
	version := ""
	verCmd := exec.Command("mysql", "--version")
	if verOutput, err := verCmd.Output(); err == nil {
		version = strings.TrimSpace(string(verOutput))
	}

	return map[string]string{
		"status":  status,
		"version": version,
	}, nil
}

func (m *MySQL) Stop() error {
	cmd := exec.Command("sudo", "systemctl", "stop", "mysql")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("停止服务失败: %v\n%s", err, string(output))
	}
	return nil
}

func (m *MySQL) GetInfo() SoftwareInfo {
	return m.info
}
