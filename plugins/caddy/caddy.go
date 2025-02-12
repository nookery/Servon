package caddy

import (
	"fmt"
	"os/exec"
	"servon/core"
	"servon/core/contract"
	"strings"
)

type Caddy struct {
	CaddyConfig
	CaddyTemplate
	*core.Core
	info contract.SoftwareInfo
}

func (c *Caddy) Install(logChan chan<- string) error {
	osType := c.GetOSType()

	switch osType {
	case core.Ubuntu, core.Debian:
		c.PrintInfo("使用 APT 包管理器安装...")
		c.PrintInfo("添加 Caddy 官方源...")

		// 下载并安装 GPG 密钥
		if err := c.RunShell("sh", "-c", "curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg"); err != nil {
			return fmt.Errorf("%s", err.Error())
		}

		// 添加 Caddy 软件源
		err := c.RunShell("sh", "-c", "curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list")
		if err != nil {
			errMsg := fmt.Sprintf("添加源失败:\n%s", err.Error())
			c.PrintErrorf("%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

		// 更新软件包索引
		if err := c.AptUpdate(); err != nil {
			errMsg := fmt.Sprintf("更新软件包索引失败: %v", err)
			c.PrintErrorf("%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

		// 安装 Caddy
		if err := c.AptInstall("caddy"); err != nil {
			errMsg := fmt.Sprintf("安装 Caddy 失败: %v", err)
			c.PrintErrorf("%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

	case core.CentOS, core.RedHat:
		errMsg := "暂不支持在 RHEL 系统上安装 Caddy"
		c.PrintErrorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)

	default:
		errMsg := fmt.Sprintf("不支持的操作系统: %s", osType)
		c.PrintErrorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 验证安装结果
	if !c.IsInstalled("caddy") {
		errMsg := "Caddy 安装验证失败，未检测到已安装的包"
		c.PrintErrorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	c.PrintSuccess("Caddy 安装完成")

	return nil
}

func (c *Caddy) Uninstall(logChan chan<- string) error {
	outputChan := make(chan string, 100)

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
		if err := c.AptRemove("caddy"); err != nil {
			return
		}

		// 清理配置文件
		if err := c.AptPurge("caddy"); err != nil {
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
	return c.StreamCommand(cmd)
}

func (c *Caddy) GetInfo() contract.SoftwareInfo {
	return c.info
}

// Reload reloads the Caddy configuration
func (c *Caddy) Reload() error {
	// 输出配置文件的存储目录
	c.PrintInfof("配置文件的存储目录: %s", c.GetConfigDir())

	// 检查 caddy 是否在运行
	running, err := c.isRunning()
	if err != nil {
		return c.PrintAndReturnError("检查 caddy 运行状态失败")
	}

	if !running {
		return c.PrintAndReturnError("Caddy 服务未运行，请先启动 Caddy")
	}

	cmd := exec.Command("caddy", "reload", "--config", c.GetCaddyfilePath())
	return c.StreamCommand(cmd)
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
func (c *Caddy) Start(logChan chan<- string) error {
	// 检查是否已安装
	if !c.IsInstalled("caddy") {
		errMsg := "Caddy 未安装，请先安装"
		c.PrintErrorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 获取当前状态
	status, err := c.GetStatus()
	if err != nil {
		errMsg := fmt.Sprintf("获取状态失败: %v", err)
		c.PrintErrorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 如果已经在运行，则不需要启动
	if status["status"] == "running" {
		c.PrintInfof("Caddy 服务已在运行中")
		return nil
	}

	c.PrintInfof("正在启动 Caddy 服务...")

	// 确保配置目录和 Caddyfile 存在
	if err := c.EnsureConfigDir(); err != nil {
		errMsg := fmt.Sprintf("创建 Caddy 配置目录失败: %v", err)
		c.PrintErrorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	if err := c.EnsureCaddyfile(); err != nil {
		errMsg := fmt.Sprintf("确保 Caddyfile 存在失败: %v", err)
		c.PrintErrorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 使用 StreamCommand 来启动 Caddy
	cmd := exec.Command("caddy", "start", "--config", c.GetCaddyfilePath())
	if err := c.StreamCommand(cmd); err != nil {
		errMsg := fmt.Sprintf("启动 Caddy 失败: %v", err)
		c.PrintErrorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	c.PrintInfof("Caddy 服务已成功启动")
	return nil
}
