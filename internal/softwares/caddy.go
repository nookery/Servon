package softwares

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type Caddy struct {
	info SoftwareInfo
}

// Configuration related constants and types
const caddyConfigTemplate = `
{{ .Domain }} {
	{{ if eq .Type "static" }}
	root * {{ .OutputPath }}
	file_server
	{{ else }}
	reverse_proxy localhost:{{ .Port }}
	{{ end }}
}
`

type Project struct {
	ID        int
	Domain    string
	Type      string
	OutputDir string
	Port      int
}

func NewCaddy() *Caddy {
	return &Caddy{
		info: SoftwareInfo{
			Name:        "caddy",
			Description: "现代化的 Web 服务器，支持自动 HTTPS",
		},
	}
}

func (c *Caddy) Install() (chan string, error) {
	outputChan := make(chan string, 100)
	apt := NewApt(outputChan)

	go func() {
		defer close(outputChan)

		// 添加 Caddy 官方源
		outputChan <- "添加 Caddy 官方源..."
		addKeyCmd := exec.Command("sudo", "apt", "install", "-y", "debian-keyring", "debian-archive-keyring", "apt-transport-https")
		if output, err := addKeyCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("添加密钥失败:\n%s", string(output))
			return
		}

		// 使用 shell 执行带管道的命令
		curlCmd := exec.Command("sh", "-c", "curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg")
		if output, err := curlCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("下载 GPG 密钥失败:\n%s", string(output))
			return
		}

		sourceCmd := exec.Command("sh", "-c", "curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list")
		if output, err := sourceCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("添加源失败:\n%s", string(output))
			return
		}

		// 更新软件包索引
		if err := apt.Update(); err != nil {
			return
		}

		// 安装 Caddy
		if err := apt.Install("caddy"); err != nil {
			return
		}

		outputChan <- "Caddy 安装完成"
	}()

	return outputChan, nil
}

func (c *Caddy) Uninstall() (chan string, error) {
	outputChan := make(chan string, 100)
	apt := NewApt(outputChan)

	go func() {
		defer close(outputChan)

		// 停止服务
		outputChan <- "停止 Caddy 服务..."
		stopCmd := exec.Command("sudo", "systemctl", "stop", "caddy")
		output, err := stopCmd.CombinedOutput()
		if err != nil {
			outputChan <- fmt.Sprintf("停止服务失败:\n%s", string(output))
		}

		// 卸载软件包及其依赖
		if err := apt.Remove("caddy"); err != nil {
			return
		}

		// 清理配置文件
		if err := apt.Purge("caddy"); err != nil {
			return
		}

		// 删除源文件
		rmSourceCmd := exec.Command("sudo", "rm", "/etc/apt/sources.list.d/caddy-stable.list")
		if output, err := rmSourceCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("删除源文件失败:\n%s", string(output))
			return
		}

		// 清理自动安装的依赖
		cleanCmd := exec.Command("sudo", "apt-get", "autoremove", "-y")
		if output, err := cleanCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("清理依赖失败:\n%s", string(output))
			return
		}

		outputChan <- "Caddy 卸载完成"
	}()

	return outputChan, nil
}

func (c *Caddy) GetStatus() (map[string]string, error) {
	dpkg := NewDpkg(nil)

	// 检查是否安装
	if !dpkg.IsInstalled("caddy") {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	// 检查服务状态
	cmd := exec.Command("systemctl", "is-active", "caddy")
	status := "stopped"
	if err := cmd.Run(); err == nil {
		status = "running"
	}

	// 获取版本
	version := ""
	verCmd := exec.Command("caddy", "version")
	if verOutput, err := verCmd.CombinedOutput(); err == nil {
		version = strings.TrimSpace(string(verOutput))
	}

	return map[string]string{
		"status":  status,
		"version": version,
	}, nil
}

func (c *Caddy) Stop() error {
	cmd := exec.Command("sudo", "systemctl", "stop", "caddy")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("停止服务失败: %v\n%s", err, string(output))
	}
	return nil
}

func (c *Caddy) GetInfo() SoftwareInfo {
	return c.info
}

// UpdateConfig updates the Caddy configuration for a project
func (c *Caddy) UpdateConfig(project *Project) error {
	// 创建配置目录
	configDir := filepath.Join("data", "caddy", "conf.d")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// 准备模板数据
	data := struct {
		Domain     string
		Type       string
		OutputPath string
		Port       int
	}{
		Domain:     project.Domain,
		Type:       project.Type,
		OutputPath: filepath.Join("data", "projects", fmt.Sprintf("%d", project.ID), project.OutputDir),
		Port:       project.Port,
	}

	// 解析并执行模板
	tmpl, err := template.New("caddy").Parse(caddyConfigTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse config template: %v", err)
	}

	// 创建配置文件
	configFile := filepath.Join(configDir, fmt.Sprintf("%d.conf", project.ID))
	f, err := os.Create(configFile)
	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to generate config file: %v", err)
	}

	// 重载配置
	return c.Reload()
}

// Reload reloads the Caddy configuration
func (c *Caddy) Reload() error {
	cmd := exec.Command("sudo", "systemctl", "reload", "caddy")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reload Caddy: %v\n%s", err, string(output))
	}
	return nil
}
