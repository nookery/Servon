package softwares

import (
	"fmt"
	"os/exec"
	"strings"
)

type Redis struct {
	info SoftwareInfo
}

func NewRedis() *Redis {
	return &Redis{
		info: SoftwareInfo{
			Name:        "redis",
			Description: "内存数据结构存储系统",
		},
	}
}

func (r *Redis) Install() (chan string, error) {
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

		// 安装 Redis
		outputChan <- "安装 Redis..."
		installCmd := exec.Command("sudo", "apt-get", "install", "-y", "redis-server")
		if output, err := installCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("安装失败:\n%s", string(output))
			return
		}

		// 启动服务
		outputChan <- "启动 Redis 服务..."
		startCmd := exec.Command("sudo", "systemctl", "start", "redis-server")
		if output, err := startCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("启动服务失败:\n%s", string(output))
			return
		}

		outputChan <- "Redis 安装完成"
	}()

	return outputChan, nil
}

func (r *Redis) Uninstall() (chan string, error) {
	outputChan := make(chan string, 100)

	go func() {
		defer close(outputChan)

		// 停止服务
		outputChan <- "停止 Redis 服务..."
		stopCmd := exec.Command("sudo", "systemctl", "stop", "redis-server")
		if output, err := stopCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("停止服务失败:\n%s", string(output))
		}

		// 卸载软件
		outputChan <- "卸载 Redis..."
		removeCmd := exec.Command("sudo", "apt-get", "remove", "-y", "redis-server")
		if output, err := removeCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("卸载失败:\n%s", string(output))
			return
		}

		// 清理配置文件
		outputChan <- "清理配置文件..."
		purgeCmd := exec.Command("sudo", "apt-get", "purge", "-y", "redis-server")
		if output, err := purgeCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("清理失败:\n%s", string(output))
			return
		}

		outputChan <- "Redis 卸载完成"
	}()

	return outputChan, nil
}

func (r *Redis) GetStatus() (map[string]string, error) {
	// 检查是否安装
	cmd := exec.Command("dpkg", "-l", "redis-server")
	if err := cmd.Run(); err != nil {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	// 检查服务状态
	statusCmd := exec.Command("systemctl", "status", "redis-server")
	output, err := statusCmd.CombinedOutput()
	status := "stopped"
	if err == nil && strings.Contains(string(output), "Active: active (running)") {
		status = "running"
	}

	// 获取版本
	version := ""
	verCmd := exec.Command("redis-server", "--version")
	if verOutput, err := verCmd.Output(); err == nil {
		version = strings.TrimSpace(string(verOutput))
	}

	return map[string]string{
		"status":  status,
		"version": version,
	}, nil
}

func (r *Redis) Stop() error {
	cmd := exec.Command("sudo", "systemctl", "stop", "redis-server")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("停止服务失败: %v\n%s", err, string(output))
	}
	return nil
}

func (r *Redis) GetInfo() SoftwareInfo {
	return r.info
}
