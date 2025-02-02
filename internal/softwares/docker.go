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
	apt := NewApt(outputChan)

	go func() {
		defer close(outputChan)

		// 更新软件包索引
		if err := apt.Update(); err != nil {
			return
		}

		// 安装必要的依赖
		if err := apt.Install("ca-certificates", "curl", "gnupg"); err != nil {
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
		if err := apt.Update(); err != nil {
			return
		}

		// 安装 Docker
		if err := apt.Install("docker-ce", "docker-ce-cli", "containerd.io", "docker-buildx-plugin", "docker-compose-plugin"); err != nil {
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
	apt := NewApt(outputChan)

	go func() {
		defer close(outputChan)

		// 停止服务
		outputChan <- "停止 Docker 服务..."
		stopCmd := exec.Command("sudo", "service", "docker", "stop")
		if output, err := stopCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("停止服务失败:\n%s", string(output))
		}

		// 卸载 Docker 包
		if err := apt.Remove("docker-ce", "docker-ce-cli", "containerd.io", "docker-buildx-plugin", "docker-compose-plugin"); err != nil {
			return
		}

		// 清理配置文件
		if err := apt.Purge("docker-ce", "docker-ce-cli", "containerd.io", "docker-buildx-plugin", "docker-compose-plugin"); err != nil {
			return
		}

		// 清理数据目录
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
	dpkg := NewDpkg(nil)

	// 检查是否安装
	if !dpkg.IsInstalled("docker-ce") {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	// 检查服务状态
	cmd := exec.Command("service", "docker", "status")
	status := "stopped"
	if err := cmd.Run(); err == nil {
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
