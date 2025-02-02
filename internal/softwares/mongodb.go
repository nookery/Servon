package softwares

import (
	"fmt"
	"os/exec"
	"strings"
)

type MongoDB struct {
	info SoftwareInfo
}

func NewMongoDB() *MongoDB {
	return &MongoDB{
		info: SoftwareInfo{
			Name:        "mongodb",
			Description: "流行的 NoSQL 数据库",
		},
	}
}

func (m *MongoDB) Install() (chan string, error) {
	outputChan := make(chan string, 100)

	go func() {
		defer close(outputChan)

		// 添加 MongoDB GPG key
		outputChan <- "添加 MongoDB GPG key..."
		keyCmd := exec.Command("sudo", "apt-key", "adv", "--keyserver", "hkp://keyserver.ubuntu.com:80", "--recv", "656408E390CFB1F5")
		if output, err := keyCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("添加 GPG key 失败:\n%s", string(output))
			return
		}

		// 添加 MongoDB 源
		outputChan <- "添加 MongoDB 源..."
		sourceCmd := exec.Command("sudo", "sh", "-c", `echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu $(lsb_release -cs)/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list`)
		if output, err := sourceCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("添加源失败:\n%s", string(output))
			return
		}

		// 更新软件包索引
		outputChan <- "更新软件包索引..."
		updateCmd := exec.Command("sudo", "apt-get", "update")
		if output, err := updateCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("更新索引失败:\n%s", string(output))
			return
		}

		// 安装 MongoDB
		outputChan <- "安装 MongoDB..."
		installCmd := exec.Command("sudo", "apt-get", "install", "-y", "mongodb-org")
		if output, err := installCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("安装失败:\n%s", string(output))
			return
		}

		outputChan <- "MongoDB 安装完成"
	}()

	return outputChan, nil
}

func (m *MongoDB) Uninstall() (chan string, error) {
	outputChan := make(chan string, 100)

	go func() {
		defer close(outputChan)

		// 停止服务
		outputChan <- "停止 MongoDB 服务..."
		stopCmd := exec.Command("sudo", "systemctl", "stop", "mongod")
		if output, err := stopCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("停止服务失败:\n%s", string(output))
		}

		// 卸载软件
		outputChan <- "卸载 MongoDB..."
		removeCmd := exec.Command("sudo", "apt-get", "remove", "-y", "mongodb-org*")
		if output, err := removeCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("卸载失败:\n%s", string(output))
			return
		}

		// 清理配置文件
		outputChan <- "清理配置文件..."
		purgeCmd := exec.Command("sudo", "apt-get", "purge", "-y", "mongodb-org*")
		if output, err := purgeCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("清理失败:\n%s", string(output))
			return
		}

		outputChan <- "MongoDB 卸载完成"
	}()

	return outputChan, nil
}

func (m *MongoDB) GetStatus() (map[string]string, error) {
	// 检查是否安装
	cmd := exec.Command("dpkg", "-l", "mongodb-org")
	if err := cmd.Run(); err != nil {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	// 检查服务状态
	statusCmd := exec.Command("systemctl", "status", "mongod")
	output, err := statusCmd.CombinedOutput()
	status := "stopped"
	if err == nil && strings.Contains(string(output), "Active: active (running)") {
		status = "running"
	}

	// 获取版本
	version := ""
	verCmd := exec.Command("mongod", "--version")
	if verOutput, err := verCmd.Output(); err == nil {
		version = strings.Split(string(verOutput), "\n")[0]
	}

	return map[string]string{
		"status":  status,
		"version": version,
	}, nil
}

func (m *MongoDB) Stop() error {
	cmd := exec.Command("sudo", "systemctl", "stop", "mongod")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("停止服务失败: %v\n%s", err, string(output))
	}
	return nil
}

func (m *MongoDB) GetInfo() SoftwareInfo {
	return m.info
}
