package software

import (
	"fmt"
	"os/exec"
	"strings"

	"servon/cmd/utils/logger"
)

type Pnpm struct {
	info SoftwareInfo
}

func NewPnpm() *Pnpm {
	return &Pnpm{
		info: SoftwareInfo{
			Name:        "pnpm",
			Description: "快速的、节省磁盘空间的包管理器",
		},
	}
}

func (p *Pnpm) Install(logChan chan<- string) error {
	logger.InfoChan(logChan, "正在安装 pnpm...")

	// 检查 nodejs 是否已安装
	nodeCmd := exec.Command("node", "--version")
	if err := nodeCmd.Run(); err != nil {
		errMsg := "请先安装 NodeJS"
		logger.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 使用 StreamCommand 来执行安装并输出详细日志
	cmd := exec.Command("npm", "install", "-g", "pnpm")
	if err := logger.StreamCommand(cmd); err != nil {
		errMsg := fmt.Sprintf("安装 pnpm 失败: %v", err)
		logger.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	logger.InfoChan(logChan, "pnpm 安装完成")
	return nil
}

func (p *Pnpm) Uninstall(logChan chan<- string) error {
	logger.InfoChan(logChan, "正在卸载 pnpm...")

	cmd := exec.Command("npm", "uninstall", "-g", "pnpm")
	if err := logger.StreamCommand(cmd); err != nil {
		errMsg := fmt.Sprintf("卸载 pnpm 失败: %v", err)
		logger.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	logger.InfoChan(logChan, "pnpm 卸载完成")
	return nil
}

func (p *Pnpm) GetStatus() (map[string]string, error) {
	// 检查 nodejs 是否已安装
	nodeCmd := exec.Command("node", "--version")
	if err := nodeCmd.Run(); err != nil {
		return map[string]string{
			"status":  "nodejs_not_installed",
			"version": "",
		}, nil
	}

	// 获取 pnpm 版本
	version := ""
	verCmd := exec.Command("pnpm", "--version")
	if output, err := verCmd.CombinedOutput(); err == nil {
		version = strings.TrimSpace(string(output))
		return map[string]string{
			"status":  "installed",
			"version": version,
		}, nil
	}

	return map[string]string{
		"status":  "not_installed",
		"version": "",
	}, nil
}

func (p *Pnpm) GetInfo() SoftwareInfo {
	return p.info
}

func (p *Pnpm) Start(logChan chan<- string) error {
	logger.InfoChan(logChan, "pnpm 是包管理工具，无需启动服务")
	return nil
}

func (p *Pnpm) Stop() error {
	logger.Info("pnpm 是包管理工具，无需停止服务")
	return nil
}
