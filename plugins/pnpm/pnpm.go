package pnpm

import (
	"fmt"
	"os/exec"
	"servon/core"
	"servon/core/contract"
	"strings"
)

func Setup(core *core.Core) {
	core.RegisterSoftware("pnpm", NewPnpm())
}

// Pnpm 实现 Software 接口
type Pnpm struct {
	info contract.SoftwareInfo
	core *core.Core
}

func NewPnpm() contract.SuperSoft {
	return &Pnpm{
		info: contract.SoftwareInfo{
			Name:        "pnpm",
			Description: "快速的、节省磁盘空间的包管理器",
		},
	}
}

// 从原来的 pnpm.go 复制所有方法实现...
func (p *Pnpm) Install(logChan chan<- string) error {
	p.core.InfoChan(logChan, "正在安装 pnpm...")

	// 检查 nodejs 是否已安装
	nodeCmd := exec.Command("node", "--version")
	if err := nodeCmd.Run(); err != nil {
		errMsg := "请先安装 NodeJS"
		p.core.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 检查 npm 是否已安装
	npmCmd := exec.Command("npm", "--version")
	if err := npmCmd.Run(); err != nil {
		errMsg := "请先安装 npm"
		p.core.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 使用 StreamCommand 来执行安装并输出详细日志
	cmd := exec.Command("npm", "install", "-g", "pnpm")
	if err := p.core.StreamCommand(cmd); err != nil {
		errMsg := fmt.Sprintf("安装 pnpm 失败: %v", err)
		p.core.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	p.core.InfoChan(logChan, "pnpm 安装完成")
	return nil
}

func (p *Pnpm) Uninstall(logChan chan<- string) error {
	p.core.InfoChan(logChan, "正在卸载 pnpm...")

	cmd := exec.Command("npm", "uninstall", "-g", "pnpm")
	if err := p.core.StreamCommand(cmd); err != nil {
		errMsg := fmt.Sprintf("卸载 pnpm 失败: %v", err)
		p.core.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	p.core.InfoChan(logChan, "pnpm 卸载完成")
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

	// 检查 npm 是否已安装
	npmCmd := exec.Command("npm", "--version")
	if err := npmCmd.Run(); err != nil {
		return map[string]string{
			"status":  "npm_not_installed",
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

func (p *Pnpm) GetInfo() contract.SoftwareInfo {
	return p.info
}

func (p *Pnpm) Start(logChan chan<- string) error {
	p.core.InfoChan(logChan, "pnpm 是包管理工具，无需启动服务")
	return nil
}

func (p *Pnpm) Stop() error {
	p.core.Info("pnpm 是包管理工具，无需停止服务")
	return nil
}
