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
	apt := NewApt(outputChan)

	go func() {
		defer close(outputChan)

		// 更新软件包索引
		if err := apt.Update(); err != nil {
			return
		}

		// 安装 Redis
		if err := apt.Install("redis-server"); err != nil {
			return
		}

		// 启动服务
		outputChan <- "启动 Redis 服务..."
		startCmd := exec.Command("sudo", "service", "redis-server", "start")
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
	apt := NewApt(outputChan)

	go func() {
		defer close(outputChan)

		// 停止服务
		outputChan <- "停止 Redis 服务..."
		stopCmd := exec.Command("sudo", "service", "redis-server", "stop")
		if output, err := stopCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("停止服务失败:\n%s", string(output))
		}

		// 卸载软件
		if err := apt.Remove("redis-server*"); err != nil {
			return
		}

		// 清理配置文件
		if err := apt.Purge("redis-server*"); err != nil {
			return
		}

		outputChan <- "Redis 卸载完成"
	}()

	return outputChan, nil
}

func (r *Redis) GetStatus() (map[string]string, error) {
	dpkg := NewDpkg(nil)

	// 检查是否安装
	if !dpkg.IsInstalled("redis-server") {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	// 检查服务状态
	cmd := exec.Command("service", "redis-server", "status")
	status := "stopped"
	if err := cmd.Run(); err == nil {
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
