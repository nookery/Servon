package pnpm

import (
	"fmt"
	"servon/core"
)

func Setup(app *core.App) {
	p := NewPnpm(app)
	app.RegisterSoftware("pnpm", p)
}

// Pnpm 实现 Software 接口
type Pnpm struct {
	info core.SoftwareInfo
	*core.App
}

func NewPnpm(app *core.App) core.SuperSoft {
	return &Pnpm{
		App: app,
		info: core.SoftwareInfo{
			Name:        "pnpm",
			Description: "快速的、节省磁盘空间的包管理器",
		},
	}
}

// 从原来的 pnpm.go 复制所有方法实现...
func (p *Pnpm) Install() error {
	p.Infof("正在安装 pnpm...")

	// 检查 nodejs 是否已安装
	err, _ := p.RunShell("node", "--version")
	if err != nil {
		errMsg := "请先安装 NodeJS"
		return p.LogAndReturnErrorf("%s", errMsg)
	}

	// 检查 npm 是否已安装
	err, _ = p.RunShell("npm", "--version")
	if err != nil {
		errMsg := "请先安装 npm"
		return p.LogAndReturnErrorf("%s", errMsg)
	}

	// 使用 StreamCommand 来执行安装并输出详细日志
	err, _ = p.RunShell("npm", "install", "-g", "pnpm")
	if err != nil {
		errMsg := fmt.Sprintf("安装 pnpm 失败: %v", err)
		return p.LogAndReturnErrorf("%s", errMsg)
	}

	p.Success("pnpm 安装完成")
	return nil
}

// Uninstall 卸载 pnpm
func (p *Pnpm) Uninstall() error {
	p.Infof("正在卸载 pnpm...")

	err, _ := p.RunShell("npm", "uninstall", "-g", "pnpm")
	if err != nil {
		errMsg := fmt.Sprintf("卸载 pnpm 失败: %v", err)
		return p.LogAndReturnErrorf("%s", errMsg)
	}

	p.Success("pnpm 卸载完成")
	return nil
}

func (p *Pnpm) GetStatus() (map[string]string, error) {
	// 检查 nodejs 是否已安装
	err, _ := p.RunShellWithOutput("node", "--version")
	if err != nil {
		return map[string]string{
			"status":  "nodejs_not_installed",
			"version": "",
		}, nil
	}

	// 检查 npm 是否已安装
	err, _ = p.RunShellWithOutput("npm", "--version")
	if err != nil {
		return map[string]string{
			"status":  "npm_not_installed",
			"version": "",
		}, nil
	}

	// 获取 pnpm 版本
	version := ""
	err, _ = p.RunShellWithOutput("pnpm", "--version")
	if err != nil {
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

func (p *Pnpm) GetInfo() core.SoftwareInfo {
	return p.info
}

func (p *Pnpm) Start() error {
	p.Infof("pnpm 是包管理工具，无需启动服务")
	return nil
}

func (p *Pnpm) Stop() error {
	p.Infof("pnpm 是包管理工具，无需停止服务")
	return nil
}
