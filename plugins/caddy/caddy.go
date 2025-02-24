package caddy

import (
	"fmt"
	"os"
	"os/exec"
	"servon/core"
	"strings"
)

type Caddy struct {
	BaseDir string
	CaddyTemplate
	*core.App
	info core.SoftwareInfo
}

func NewCaddy(app *core.App) *Caddy {
	return &Caddy{
		App: app,
		info: core.SoftwareInfo{
			Name:            "caddy",
			Description:     "Modern web server with automatic HTTPS",
			IsProxySoftware: true,
			IsGateway:       true,
		},
	}
}

// GetInfo 获取软件信息
func (c *Caddy) GetInfo() core.SoftwareInfo {
	return c.info
}

func (c *Caddy) Install() error {
	osType := c.GetOSType()

	switch osType {
	case core.Ubuntu, core.Debian:
		c.Info("使用 APT 包管理器安装...")
		c.Info("添加 Caddy 官方源...")

		// 下载和安装 GPG 密钥
		c.Info("下载和安装 GPG 密钥...")
		err, _ := c.RunShell("sh", "-c", "curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg")
		if err != nil {
			c.Errorf("下载 GPG 密钥失败: %v", err)
			return err
		}

		// 添加 Caddy 软件源
		c.Info("添加 Caddy 软件源...")
		err, _ = c.RunShell("sh", "-c", "curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list")
		if err != nil {
			c.Errorf("添加 Caddy 软件源失败: %v", err)
			return err
		}

		// 更新软件包索引
		c.Info("更新软件包索引...")
		if err := c.AptUpdate(); err != nil {
			c.Errorf("更新软件包索引失败: %v", err)
			return err
		}

		// 安装 Caddy
		c.Info("安装 Caddy...")
		if err := c.AptInstall("caddy"); err != nil {
			c.Errorf("Caddy 安装失败: %v", err)
			return err
		}

		c.Success("Caddy 安装完成")
	case core.CentOS, core.RedHat:
		errMsg := "暂不支持在 RHEL 系统上安装 Caddy"
		c.Errorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)

	default:
		errMsg := fmt.Sprintf("不支持的操作系统: %s", osType)
		c.Errorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 验证安装结果
	if !c.IsInstalled("caddy") {
		errMsg := "Caddy 安装验证失败，未检测到已安装的包"
		c.Errorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	c.Success("Caddy 安装完成")

	return nil
}

// Uninstall 卸载 Caddy
func (c *Caddy) Uninstall() error {
	// 停止服务
	c.Info("停止 Caddy 服务...")
	stopCmd := exec.Command("sudo", "systemctl", "stop", "caddy")
	if err, _ := c.RunShell(stopCmd.String()); err != nil {
		return fmt.Errorf("停止服务失败:\n%s", err)
	}

	// 卸载软件包及其依赖
	c.Info("卸载软件包及其依赖...")
	if err := c.AptRemove("caddy"); err != nil {
		return fmt.Errorf("卸载软件包及其依赖失败:\n%s", err)
	}

	// 清理配置文件
	if err := c.AptPurge("caddy"); err != nil {
		return fmt.Errorf("清理配置文件失败:\n%s", err)
	}

	// 删除源文件
	err, _ := c.RunShell("sudo", "rm", "/etc/apt/sources.list.d/caddy-stable.list")
	if err != nil {
		return fmt.Errorf("删除源文件失败:\n%s", err)
	}

	// 清理自动安装的依赖
	cleanCmd := exec.Command("sudo", "apt-get", "autoremove", "-y")
	if err := cleanCmd.Run(); err != nil {
		return fmt.Errorf("清理依赖失败:\n%s", err)
	}

	c.Success("Caddy 卸载完成")

	return nil
}

func (c *Caddy) GetStatus() (map[string]string, error) {
	if !c.IsInstalled("caddy") {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	status := "stopped"
	if c.IsActive("caddy") {
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
	err, _ := c.RunShell(cmd.String())
	return err
}

// Reload reloads the Caddy configuration
func (c *Caddy) Reload() error {
	// 输出配置文件的存储目录
	c.Infof("配置文件的存储目录: %s", c.GetConfigDir())

	// 检查 caddy 是否在运行
	running, err := c.isRunning()
	if err != nil {
		return c.LogAndReturnError("检查 caddy 运行状态失败")
	}

	if !running {
		return c.LogAndReturnError("Caddy 服务未运行，请先启动 Caddy")
	}

	cmd := exec.Command("caddy", "reload", "--config", c.GetCaddyfilePath())
	err, _ = c.RunShell(cmd.String())
	return err
}

// isRunning 检查 caddy 是否在运行
func (c *Caddy) isRunning() (bool, error) {
	cmd := exec.Command("pgrep", "caddy")
	if err := cmd.Run(); err != nil {
		// exit status 1 表示没有找到进程
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Start starts the Caddy service
func (c *Caddy) Start() error {
	// 检查是否已安装
	if !c.IsInstalled("caddy") {
		errMsg := "Caddy 未安装，请先安装"
		c.Errorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 获取当前状态
	status, err := c.GetStatus()
	if err != nil {
		errMsg := fmt.Sprintf("获取状态失败: %v", err)
		c.Errorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 如果已经在运行，则不需要启动
	if status["status"] == "running" {
		c.Infof("Caddy 服务已在运行中")
		return nil
	}

	c.Infof("正在启动 Caddy 服务...")

	// 确保配置目录和 Caddyfile 存在
	if err := c.EnsureConfigDir(); err != nil {
		errMsg := fmt.Sprintf("创建 Caddy 配置目录失败: %v", err)
		return c.LogAndReturnErrorf("%s", errMsg)
	}

	if err := c.EnsureCaddyfile(); err != nil {
		errMsg := fmt.Sprintf("确保 Caddyfile 存在失败: %v", err)
		return c.LogAndReturnErrorf("%s", errMsg)
	}

	// 使用 StreamCommand 来启动 Caddy
	cmd := exec.Command("caddy", "start", "--config", c.GetCaddyfilePath())
	err, _ = c.RunShell(cmd.String())
	if err != nil {
		errMsg := fmt.Sprintf("启动 Caddy 失败: %v", err)
		return c.LogAndReturnErrorf("%s", errMsg)
	}

	c.Infof("Caddy 服务已成功启动")
	return nil
}

// Gateway 接口实现
func (c *Caddy) GetConfig() (map[string]interface{}, error) {
	// 确保配置目录和文件存在
	if err := c.EnsureConfigDir(); err != nil {
		return nil, fmt.Errorf("确保配置目录存在失败: %v", err)
	}

	if err := c.EnsureCaddyfile(); err != nil {
		return nil, fmt.Errorf("确保配置文件存在失败: %v", err)
	}

	content, err := os.ReadFile(c.GetCaddyfilePath())
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"config": string(content),
	}, nil
}

func (c *Caddy) SetConfig(config map[string]interface{}) error {
	if configStr, ok := config["config"].(string); ok {
		err := os.WriteFile(c.GetCaddyfilePath(), []byte(configStr), 0644)
		if err != nil {
			return err
		}
		return c.Reload()
	}
	return fmt.Errorf("invalid config format")
}

func (c *Caddy) GetProjects() ([]core.Project, error) {
	// TODO: 解析 Caddyfile 获取项目列表
	return []core.Project{}, nil
}

func (c *Caddy) AddProject(project core.Project) error {
	// TODO: 将项目配置添加到 Caddyfile
	return c.Reload()
}

func (c *Caddy) RemoveProject(projectName string) error {
	// TODO: 从 Caddyfile 中移除项目配置
	return c.Reload()
}

func (c *Caddy) ReloadConfig() error {
	return c.Reload()
}
