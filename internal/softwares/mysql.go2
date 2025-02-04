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
	apt := NewApt(outputChan)

	go func() {
		defer close(outputChan)

		// 更新软件包索引
		if err := apt.Update(); err != nil {
			return
		}

		// 安装 MySQL
		if err := apt.Install("mysql-server"); err != nil {
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
	apt := NewApt(outputChan)

	go func() {
		defer close(outputChan)

		// 停止服务
		outputChan <- "停止 MySQL 服务..."
		stopCmd := exec.Command("sudo", "service", "mysql", "stop")
		if output, err := stopCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("停止服务失败:\n%s", string(output))
		}

		// 卸载软件
		if err := apt.Remove("mysql-server*"); err != nil {
			return
		}

		// 清理配置文件
		if err := apt.Purge("mysql-server*"); err != nil {
			return
		}

		outputChan <- "MySQL 卸载完成"
	}()

	return outputChan, nil
}

func (m *MySQL) GetStatus() (map[string]string, error) {
	dpkg := NewDpkg(nil)

	// 检查是否安装
	if !dpkg.IsInstalled("mysql-server") {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	// 检查服务状态
	cmd := exec.Command("service", "mysql", "status")
	status := "stopped"
	if err := cmd.Run(); err == nil {
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
	cmd := exec.Command("sudo", "service", "mysql", "stop")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("停止服务失败: %v\n%s", err, string(output))
	}
	return nil
}

func (m *MySQL) GetInfo() SoftwareInfo {
	return m.info
}
