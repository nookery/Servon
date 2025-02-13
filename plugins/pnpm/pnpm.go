package pnpm

import (
	"fmt"
	"servon/core"
	"servon/core/contract"
	"strings"
)

func Setup(core *core.Core) {
	p := NewPnpm(core)
	core.RegisterSoftware("pnpm", p)
}

// Pnpm 实现 Software 接口
type Pnpm struct {
	info contract.SoftwareInfo
	*core.Core
}

func NewPnpm(core *core.Core) contract.SuperSoft {
	return &Pnpm{
		Core: core,
		info: contract.SoftwareInfo{
			Name:        "pnpm",
			Description: "快速的、节省磁盘空间的包管理器",
		},
	}
}

// 从原来的 pnpm.go 复制所有方法实现...
func (p *Pnpm) Install() error {
	p.PrintInfo("正在安装 pnpm...")

	// 检查 nodejs 是否已安装
	if err := p.RunShell("node", "--version"); err != nil {
		errMsg := "请先安装 NodeJS"
		p.PrintErrorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 检查 npm 是否已安装
	if err := p.RunShell("npm", "--version"); err != nil {
		errMsg := "请先安装 npm"
		p.PrintErrorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 使用 StreamCommand 来执行安装并输出详细日志
	if err := p.RunShell("npm", "install", "-g", "pnpm"); err != nil {
		errMsg := fmt.Sprintf("安装 pnpm 失败: %v", err)
		p.PrintErrorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	p.PrintSuccess("pnpm 安装完成")
	return nil
}

// Uninstall 卸载 pnpm
func (p *Pnpm) Uninstall() error {
	p.PrintInfo("正在卸载 pnpm...")

	if err := p.RunShell("npm", "uninstall", "-g", "pnpm"); err != nil {
		errMsg := fmt.Sprintf("卸载 pnpm 失败: %v", err)
		p.PrintErrorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	p.PrintSuccess("pnpm 卸载完成")
	return nil
}

func (p *Pnpm) GetStatus() (map[string]string, error) {
	// 检查 nodejs 是否已安装
	if _, err := p.RunShellWithOutput("node", "--version"); err != nil {
		return map[string]string{
			"status":  "nodejs_not_installed",
			"version": "",
		}, nil
	}

	// 检查 npm 是否已安装
	if _, err := p.RunShellWithOutput("npm", "--version"); err != nil {
		return map[string]string{
			"status":  "npm_not_installed",
			"version": "",
		}, nil
	}

	// 获取 pnpm 版本
	version := ""
	if output, err := p.RunShellWithOutput("pnpm", "--version"); err == nil {
		version = strings.TrimSpace(output)
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

func (p *Pnpm) Start() error {
	p.PrintInfo("pnpm 是包管理工具，无需启动服务")
	return nil
}

func (p *Pnpm) Stop() error {
	p.PrintInfo("pnpm 是包管理工具，无需停止服务")
	return nil
}
