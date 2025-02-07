package softwares

import (
	"fmt"
	"os/exec"
	"strings"

	"servon/cmd/internal/system"
	"servon/cmd/internal/utils"
)

type Caddy struct {
	info   SoftwareInfo
	config *CaddyConfig
}

type Project struct {
	Name      string
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
		config: NewCaddyConfig(),
	}
}

func (c *Caddy) Install(logChan chan<- string) error {
	outputChan := make(chan string, 100)
	apt := system.NewApt()

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
	dpkg := system.NewDpkg(outputChan)
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
	apt := system.NewApt()

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
	dpkg := system.NewDpkg(nil)

	if !dpkg.IsInstalled("caddy") {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	status := "stopped"
	if system.ServiceIsActive("caddy") {
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
	cmd := exec.Command("caddy", "stop")
	return utils.StreamCommand(cmd)
}

func (c *Caddy) GetInfo() SoftwareInfo {
	return c.info
}

// GetConfigDir returns the base configuration directory for Caddy
func (c *Caddy) GetConfigDir() string {
	return c.config.GetConfigDir()
}

// GetCaddyfilePath returns the path to the main Caddyfile
func (c *Caddy) GetCaddyfilePath() string {
	return c.config.GetCaddyfilePath()
}

// GetProjectConfigPath returns the configuration file path for a specific project
func (c *Caddy) GetProjectConfigPath(projectName string) string {
	return c.config.GetProjectConfigPath(projectName)
}

// EnsureConfigDir ensures the Caddy configuration directory exists
func (c *Caddy) EnsureConfigDir() error {
	return c.config.EnsureConfigDir()
}

// EnsureCaddyfile ensures the main Caddyfile exists, creating it from template if needed
func (c *Caddy) EnsureCaddyfile() error {
	return c.config.EnsureCaddyfile()
}

// UpdateConfig updates the Caddy configuration for a project
func (c *Caddy) UpdateConfig(project *Project) error {
	if err := c.config.UpdateProjectConfig(project); err != nil {
		return err
	}
	return c.Reload()
}

// Reload reloads the Caddy configuration
func (c *Caddy) Reload() error {
	cmd := exec.Command("caddy", "reload", "--config", c.config.GetCaddyfilePath())
	return utils.StreamCommand(cmd)
}

// Start starts the Caddy service
func (c *Caddy) Start(logChan chan<- string) error {
	// 检查是否已安装
	dpkg := system.NewDpkg(nil)
	if !dpkg.IsInstalled("caddy") {
		errMsg := "Caddy 未安装，请先安装"
		utils.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 获取当前状态
	status, err := c.GetStatus()
	if err != nil {
		errMsg := fmt.Sprintf("获取状态失败: %v", err)
		utils.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 如果已经在运行，则不需要启动
	if status["status"] == "running" {
		utils.InfoChan(logChan, "Caddy 服务已在运行中")
		return nil
	}

	utils.DebugChan(logChan, "正在启动 Caddy 服务...")

	// 确保配置目录和 Caddyfile 存在
	if err := c.config.EnsureConfigDir(); err != nil {
		errMsg := fmt.Sprintf("创建 Caddy 配置目录失败: %v", err)
		utils.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	if err := c.config.EnsureCaddyfile(); err != nil {
		errMsg := fmt.Sprintf("确保 Caddyfile 存在失败: %v", err)
		utils.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 使用 StreamCommand 来启动 Caddy
	cmd := exec.Command("caddy", "start", "--config", c.config.GetCaddyfilePath())
	if err := utils.StreamCommand(cmd); err != nil {
		errMsg := fmt.Sprintf("启动 Caddy 失败: %v", err)
		utils.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	utils.DebugChan(logChan, "Caddy 服务已成功启动")
	return nil
}
