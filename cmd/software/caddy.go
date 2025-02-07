package software

import (
	"fmt"
	"os/exec"
	"strings"

	"servon/cmd/contract"
	"servon/cmd/system"
	"servon/cmd/utils/logger"
)

type Caddy struct {
	info   contract.SoftwareInfo
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
		info: contract.SoftwareInfo{
			Name:        "caddy",
			Description: "现代化的 Web 服务器，支持自动 HTTPS",
		},
		config: NewCaddyConfig(),
	}
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

func (c *Caddy) GetInfo() contract.SoftwareInfo {
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
	// 检查 caddy 是否在运行
	running, err := c.isRunning()
	if err != nil {
		return fmt.Errorf("检查 caddy 运行状态失败: %v", err)
	}

	if !running {
		logger.Warn("Caddy 服务未运行，请先启动 Caddy")
		return nil
	}

	cmd := exec.Command("caddy", "reload", "--config", c.config.GetCaddyfilePath())
	return logger.StreamCommand(cmd)
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
