package softwares

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"servon/internal/system"
	"servon/internal/utils"
)

type Clash struct {
	info SoftwareInfo
}

// Configuration related constants and types
const clashConfigTemplate = `
port: 7890
socks-port: 7891
allow-lan: true
mode: Rule
log-level: info
external-controller: :9090
proxies:
  # Configure your proxy servers here
proxy-groups:
  # Configure your proxy groups here
rules:
  # Configure your rules here
`

func NewClash() *Clash {
	return &Clash{
		info: SoftwareInfo{
			Name:        "clash",
			Description: "A rule-based tunnel in Go",
		},
	}
}

func (c *Clash) Install(logChan chan<- string) error {
	osType := utils.GetOSType()
	utils.InfoChan(logChan, "检测到操作系统: %s", osType)

	switch osType {
	case utils.Ubuntu, utils.Debian:
		utils.InfoChan(logChan, "开始安装 Clash...")

		// 创建安装目录
		installDir := "/usr/local/bin"
		if err := os.MkdirAll(installDir, 0755); err != nil {
			errMsg := fmt.Sprintf("创建安装目录失败: %v", err)
			utils.ErrorChan(logChan, "%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

		// 下载最新版本的 Clash
		downloadCmd := exec.Command("curl", "-L",
			"https://github.com/Dreamacro/clash/releases/download/v1.18.0/clash-linux-amd64-v1.18.0.gz",
			"-o", "/tmp/clash.gz")
		if err := utils.StreamCommand(downloadCmd); err != nil {
			return fmt.Errorf("%s", err)
		}

		// 解压
		gunzipCmd := exec.Command("gunzip", "-f", "/tmp/clash.gz")
		if err := utils.StreamCommand(gunzipCmd); err != nil {
			return fmt.Errorf("%s", err)
		}

		// 移动到安装目录并设置权限
		moveCmd := exec.Command("sudo", "mv", "/tmp/clash", filepath.Join(installDir, "clash"))
		if err := utils.StreamCommand(moveCmd); err != nil {
			return fmt.Errorf("%s", err)
		}

		chmodCmd := exec.Command("sudo", "chmod", "+x", filepath.Join(installDir, "clash"))
		if err := utils.StreamCommand(chmodCmd); err != nil {
			return fmt.Errorf("%s", err)
		}

		// 创建系统服务
		serviceContent := `[Unit]
Description=Clash Daemon
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/clash -d /etc/clash
Restart=on-failure

[Install]
WantedBy=multi-user.target`

		if err := os.WriteFile("/etc/systemd/system/clash.service", []byte(serviceContent), 0644); err != nil {
			errMsg := fmt.Sprintf("创建服务文件失败: %v", err)
			utils.ErrorChan(logChan, "%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

		// 创建配置目录
		if err := os.MkdirAll("/etc/clash", 0755); err != nil {
			errMsg := fmt.Sprintf("创建配置目录失败: %v", err)
			utils.ErrorChan(logChan, "%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

		// 创建默认配置文件
		if err := os.WriteFile("/etc/clash/config.yaml", []byte(clashConfigTemplate), 0644); err != nil {
			errMsg := fmt.Sprintf("创建配置文件失败: %v", err)
			utils.ErrorChan(logChan, "%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

		// 重载系统服务
		reloadCmd := exec.Command("sudo", "systemctl", "daemon-reload")
		if output, err := reloadCmd.CombinedOutput(); err != nil {
			errMsg := fmt.Sprintf("重载系统服务失败:\n%s", string(output))
			utils.ErrorChan(logChan, "%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

	case utils.CentOS, utils.RedHat:
		errMsg := "暂不支持在 RHEL 系统上安装 Clash"
		utils.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)

	default:
		errMsg := fmt.Sprintf("不支持的操作系统: %s", osType)
		utils.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	utils.InfoChan(logChan, "Clash 安装完成")
	return nil
}

func (c *Clash) Uninstall(logChan chan<- string) error {
	outputChan := make(chan string, 100)

	go func() {
		defer close(outputChan)

		// 停止服务
		outputChan <- "停止 Clash 服务..."
		stopCmd := exec.Command("sudo", "systemctl", "stop", "clash")
		if output, err := stopCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("停止服务失败:\n%s", string(output))
		}

		// 禁用服务
		disableCmd := exec.Command("sudo", "systemctl", "disable", "clash")
		if output, err := disableCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("禁用服务失败:\n%s", string(output))
		}

		// 删除服务文件
		rmServiceCmd := exec.Command("sudo", "rm", "/etc/systemd/system/clash.service")
		if output, err := rmServiceCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("删除服务文件失败:\n%s", string(output))
		}

		// 删除二进制文件
		rmBinCmd := exec.Command("sudo", "rm", "/usr/local/bin/clash")
		if output, err := rmBinCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("删除二进制文件失败:\n%s", string(output))
		}

		// 删除配置目录
		rmConfigCmd := exec.Command("sudo", "rm", "-rf", "/etc/clash")
		if output, err := rmConfigCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("删除配置目录失败:\n%s", string(output))
		}

		outputChan <- "Clash 卸载完成"
	}()

	return nil
}

func (c *Clash) GetStatus() (map[string]string, error) {
	// 检查二进制文件是否存在
	if _, err := os.Stat("/usr/local/bin/clash"); os.IsNotExist(err) {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	status := "stopped"
	if system.ServiceIsActive("clash") {
		status = "running"
	}

	// 获取版本
	version := ""
	verCmd := exec.Command("clash", "-v")
	if verOutput, err := verCmd.CombinedOutput(); err == nil {
		version = strings.TrimSpace(string(verOutput))
	}

	return map[string]string{
		"status":  status,
		"version": version,
	}, nil
}

func (c *Clash) Stop() error {
	return system.ServiceStop("clash")
}

func (c *Clash) GetInfo() SoftwareInfo {
	return c.info
}

func (c *Clash) Start(logChan chan<- string) error {
	// 检查是否已安装
	if _, err := os.Stat("/usr/local/bin/clash"); os.IsNotExist(err) {
		errMsg := "Clash 未安装，请先安装"
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
		utils.InfoChan(logChan, "Clash 服务已在运行中")
		return nil
	}

	utils.DebugChan(logChan, "正在启动 Clash 服务...")

	// 启动服务
	if err := system.ServiceStart("clash"); err != nil {
		errMsg := fmt.Sprintf("启动 Clash 失败: %v", err)
		utils.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	utils.DebugChan(logChan, "Clash 服务已成功启动")
	return nil
}

func (c *Clash) Reload() error {
	return system.ServiceReload("clash")
}
