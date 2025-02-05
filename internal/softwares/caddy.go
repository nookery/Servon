package softwares

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"servon/internal/system"
	"servon/internal/utils"
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

func (c *Caddy) Install(logChan chan<- string) error {
	outputChan := make(chan string, 100)
	apt := NewApt()

	// 检查操作系统类型
	osType := utils.GetOSType()
	utils.InfoChan(logChan, "检测到操作系统: %s", osType)

	switch osType {
	case utils.Ubuntu, utils.Debian:
		utils.InfoChan(logChan, "使用 APT 包管理器安装...")
		utils.InfoChan(logChan, "添加 Caddy 官方源...")

		// 下载并安装 GPG 密钥
		curlCmd := exec.Command("sh", "-c", "curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg")
		if output, err := curlCmd.CombinedOutput(); err != nil {
			errMsg := fmt.Sprintf("下载 GPG 密钥失败:\n%s", string(output))
			utils.ErrorChan(outputChan, "%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

		// 添加 Caddy 软件源
		sourceCmd := exec.Command("sh", "-c", "curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list")
		if output, err := sourceCmd.CombinedOutput(); err != nil {
			errMsg := fmt.Sprintf("添加源失败:\n%s", string(output))
			utils.ErrorChan(outputChan, "%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

		// 更新软件包索引
		if err := apt.Update(); err != nil {
			errMsg := fmt.Sprintf("更新软件包索引失败: %v", err)
			utils.ErrorChan(outputChan, "%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

		// 安装 Caddy
		if err := apt.Install("caddy"); err != nil {
			errMsg := fmt.Sprintf("安装 Caddy 失败: %v", err)
			utils.ErrorChan(outputChan, "%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

	case utils.CentOS, utils.RedHat:
		errMsg := "暂不支持在 RHEL 系统上安装 Caddy"
		utils.ErrorChan(outputChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)

	default:
		errMsg := fmt.Sprintf("不支持的操作系统: %s", osType)
		utils.ErrorChan(outputChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 验证安装结果
	dpkg := NewDpkg(nil)
	if !dpkg.IsInstalled("caddy") {
		errMsg := "Caddy 安装验证失败，未检测到已安装的包"
		utils.ErrorChan(outputChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	utils.InfoChan(outputChan, "Caddy 安装完成")

	return nil
}

func (c *Caddy) Uninstall(logChan chan<- string) error {
	outputChan := make(chan string, 100)
	apt := NewApt()

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

	return nil
}

func (c *Caddy) GetStatus() (map[string]string, error) {
	dpkg := NewDpkg(nil)

	if !dpkg.IsInstalled("caddy") {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	serviceManager := system.NewServiceManager()
	utils.Debug("Using %s to check caddy status", serviceManager.Type())

	status := "stopped"
	if serviceManager.IsActive("caddy") {
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
	serviceManager := system.NewServiceManager()
	return serviceManager.Stop("caddy")
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
	serviceManager := system.NewServiceManager()
	return serviceManager.Reload("caddy")
}

// Start starts the Caddy service
func (c *Caddy) Start(logChan chan<- string) error {
	outputChan := make(chan string, 100)

	go func() {
		defer close(outputChan)

		// 检查是否已安装
		dpkg := NewDpkg(nil)
		if !dpkg.IsInstalled("caddy") {
			errMsg := "Caddy 未安装，请先安装"
			utils.ErrorChan(outputChan, "%s", errMsg)
			return
		}

		// 获取当前状态
		status, err := c.GetStatus()
		if err != nil {
			errMsg := fmt.Sprintf("获取状态失败: %v", err)
			utils.ErrorChan(outputChan, "%s", errMsg)
			return
		}

		// 如果已经在运行，则不需要启动
		if status["status"] == "running" {
			utils.InfoChan(outputChan, "Caddy 服务已在运行中")
			return
		}

		utils.InfoChan(outputChan, "正在启动 Caddy 服务...")

		// 启动服务
		serviceManager := system.NewServiceManager()
		if err := serviceManager.Start("caddy"); err != nil {
			errMsg := fmt.Sprintf("启动服务失败: %v", err)
			utils.ErrorChan(outputChan, "%s", errMsg)
			return
		}

		// 验证服务是否成功启动
		if !serviceManager.IsActive("caddy") {
			errMsg := "服务启动失败"
			utils.ErrorChan(outputChan, "%s", errMsg)
			return
		}

		utils.InfoChan(outputChan, "Caddy 服务已成功启动")
	}()

	return nil
}
