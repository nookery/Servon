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
	apt := NewApt(outputChan)

	go func() {
		defer close(outputChan)

		// 更新软件包索引
		if err := apt.Update(); err != nil {
			return
		}

		// 安装 Nginx
		if err := apt.Install("nginx"); err != nil {
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
	apt := NewApt(outputChan)

	go func() {
		defer close(outputChan)

		// 停止服务
		outputChan <- "停止 Nginx 服务..."
		stopCmd := exec.Command("sudo", "service", "nginx", "stop")
		output, err := stopCmd.CombinedOutput()
		if err != nil {
			outputChan <- fmt.Sprintf("停止服务失败:\n%s", string(output))
		}

		// 卸载软件包及其依赖
		if err := apt.Remove("nginx*"); err != nil {
			return
		}

		// 清理配置文件
		if err := apt.Purge("nginx*"); err != nil {
			return
		}

		// 清理自动安装的依赖
		cleanCmd := exec.Command("sudo", "apt-get", "autoremove", "-y")
		if output, err := cleanCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("清理依赖失败:\n%s", string(output))
			return
		}

		outputChan <- "Nginx 卸载完成"
	}()

	return outputChan, nil
}

func (n *Nginx) GetStatus() (map[string]string, error) {
	dpkg := NewDpkg(nil)

	// 检查是否安装
	if !dpkg.IsInstalled("nginx") {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	// 检查服务状态
	cmd := exec.Command("service", "nginx", "status")
	status := "stopped"
	if err := cmd.Run(); err == nil {
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
