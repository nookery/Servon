package software

import (
	"fmt"
	"os/exec"
	"strings"

	"servon/cmd/utils/logger"
)

type Npm struct {
	info SoftwareInfo
}

func NewNpm() *Npm {
	return &Npm{
		info: SoftwareInfo{
			Name:        "npm",
			Description: "Node.js 默认的包管理器",
		},
	}
}

func (n *Npm) Install(logChan chan<- string) error {
	logger.InfoChan(logChan, "正在检查 npm...")

	// 检查 nodejs 是否已安装
	nodeCmd := exec.Command("node", "--version")
	if err := nodeCmd.Run(); err != nil {
		errMsg := "请先安装 NodeJS"
		logger.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 检查 npm 是否已安装
	npmCmd := exec.Command("npm", "--version")
	if err := npmCmd.Run(); err != nil {
		logger.InfoChan(logChan, "npm 未安装，正在通过 apt 安装...")

		// 使用 apt 安装 npm
		cmd := exec.Command("apt", "install", "-y", "npm")
		if err := logger.StreamCommand(cmd); err != nil {
			errMsg := "npm 安装失败"
			logger.ErrorChan(logChan, "%s: %v", errMsg, err)
			return fmt.Errorf("%s: %v", errMsg, err)
		}
	}

	logger.InfoChan(logChan, "npm 已安装")
	return nil
}

func (n *Npm) Uninstall(logChan chan<- string) error {
	logger.InfoChan(logChan, "npm 是 NodeJS 的一部分，无法单独卸载")
	return nil
}

func (n *Npm) GetStatus() (map[string]string, error) {
	// 检查 nodejs 是否已安装
	nodeCmd := exec.Command("node", "--version")
	if err := nodeCmd.Run(); err != nil {
		return map[string]string{
			"status":  "nodejs_not_installed",
			"version": "",
		}, nil
	}

	// 获取 npm 版本
	version := ""
	verCmd := exec.Command("npm", "--version")
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

func (n *Npm) GetInfo() SoftwareInfo {
	return n.info
}

func (n *Npm) Start(logChan chan<- string) error {
	logger.InfoChan(logChan, "npm 是包管理工具，无需启动服务")
	return nil
}

func (n *Npm) Stop() error {
	logger.Info("npm 是包管理工具，无需停止服务")
	return nil
}
