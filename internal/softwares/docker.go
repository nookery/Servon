package softwares

import (
	"fmt"
	"os/exec"
	"strings"
)

type Docker struct {
	info SoftwareInfo
}

func NewDocker() *Docker {
	return &Docker{
		info: SoftwareInfo{
			Name:        "docker",
			Description: "应用容器引擎",
		},
	}
}

func (d *Docker) Install() (chan string, error) {
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

		// 安装必要的依赖
		outputChan <- "安装依赖包..."
		depsCmd := exec.Command("sudo", "apt-get", "install", "-y", "ca-certificates", "curl", "gnupg")
		if output, err := depsCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("安装依赖失败:\n%s", string(output))
			return
		}

		// 添加 Docker 官方 GPG 密钥
		outputChan <- "添加 Docker GPG 密钥..."
		keyCmd := exec.Command("sudo", "install", "-m", "0755", "-d", "/etc/apt/keyrings")
		keyCmd.Run() // 忽略错误，目录可能已存在

		// 下载并安装 GPG 密钥
		downloadCmd := exec.Command("sh", "-c", "curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg")
		if output, err := downloadCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("下载 GPG 密钥失败:\n%s", string(output))
			return
		}

		// 设置仓库
		outputChan <- "设置 Docker 仓库..."
		repoCmd := exec.Command("sh", "-c", `echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null`)
		if output, err := repoCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("设置仓库失败:\n%s", string(output))
			return
		}

		// 再次更新软件包索引
		outputChan <- "更新软件包索引..."
		updateCmd = exec.Command("sudo", "apt-get", "update")
		if output, err := updateCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("更新索引失败:\n%s", string(output))
			return
		}

		// 安装 Docker
		outputChan <- "安装 Docker..."
		installCmd := exec.Command("sudo", "apt-get", "install", "-y", "docker-ce", "docker-ce-cli", "containerd.io", "docker-buildx-plugin", "docker-compose-plugin")
		if output, err := installCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("安装失败:\n%s", string(output))
			return
		}

		// 将当前用户添加到 docker 组
		outputChan <- "配置用户权限..."
		userCmd := exec.Command("sudo", "usermod", "-aG", "docker", "$USER")
		if output, err := userCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("配置用户权限失败:\n%s", string(output))
			return
		}

		outputChan <- "Docker 安装完成"
	}()

	return outputChan, nil
}

func (d *Docker) Uninstall() (chan string, error) {
	outputChan := make(chan string, 100)

	go func() {
		defer close(outputChan)

		// 停止服务
		outputChan <- "停止 Docker 服务..."
		stopCmd := exec.Command("sudo", "systemctl", "stop", "docker")
		if output, err := stopCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("停止服务失败:\n%s", string(output))
		}

		// 卸载软件包
		outputChan <- "卸载 Docker..."
		removeCmd := exec.Command("sudo", "apt-get", "remove", "-y", "docker-ce", "docker-ce-cli", "containerd.io", "docker-buildx-plugin", "docker-compose-plugin")
		if output, err := removeCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("卸载失败:\n%s", string(output))
			return
		}

		// 清理配置和数据
		outputChan <- "清理 Docker 数据..."
		purgeCmd := exec.Command("sudo", "rm", "-rf", "/var/lib/docker", "/etc/docker", "/var/run/docker.sock")
		if output, err := purgeCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("清理失败:\n%s", string(output))
			return
		}

		outputChan <- "Docker 卸载完成"
	}()

	return outputChan, nil
}

func (d *Docker) GetStatus() (map[string]string, error) {
	// 检查是否安装
	cmd := exec.Command("dpkg", "-l", "docker-ce")
	if err := cmd.Run(); err != nil {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	// 检查服务状态
	statusCmd := exec.Command("systemctl", "status", "docker")
	output, err := statusCmd.CombinedOutput()
	status := "stopped"
	if err == nil && strings.Contains(string(output), "Active: active (running)") {
		status = "running"
	}

	// 获取版本
	version := ""
	verCmd := exec.Command("docker", "--version")
	if verOutput, err := verCmd.Output(); err == nil {
		version = strings.TrimSpace(string(verOutput))
	}

	return map[string]string{
		"status":  status,
		"version": version,
	}, nil
}

func (d *Docker) Stop() error {
	cmd := exec.Command("sudo", "systemctl", "stop", "docker")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("停止服务失败: %v\n%s", err, string(output))
	}
	return nil
}

func (d *Docker) GetInfo() SoftwareInfo {
	return d.info
}
