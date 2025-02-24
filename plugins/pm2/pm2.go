package pm2

import (
	"fmt"
	"servon/core"
)

func Setup(app *core.App) {
	p := NewPm2(app)
	app.RegisterSoftware("pm2", p)
}

// Pm2 实现 Software 接口
type Pm2 struct {
	info core.SoftwareInfo
	*core.App
}

func NewPm2(app *core.App) core.SuperSoft {
	return &Pm2{
		App: app,
		info: core.SoftwareInfo{
			Name:        "pnpm",
			Description: "快速的、节省磁盘空间的包管理器",
		},
	}
}

// Install 安装 pm2
func (p *Pm2) Install() error {
	p.Infof("正在安装 pm2...")

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

	err, _ = p.RunShell("npm", "install", "-g", "pm2")
	if err != nil {
		errMsg := fmt.Sprintf("安装 pm2 失败: %v", err)
		return p.LogAndReturnErrorf("%s", errMsg)
	}

	p.Success("pm2 安装完成")
	return nil
}

// Uninstall 卸载 pm2
func (p *Pm2) Uninstall() error {
	p.Infof("正在卸载 pm2...")

	err, _ := p.RunShell("npm", "uninstall", "-g", "pm2")
	if err != nil {
		errMsg := fmt.Sprintf("卸载 pm2 失败: %v", err)
		return p.LogAndReturnErrorf("%s", errMsg)
	}

	p.Success("pm2 卸载完成")
	return nil
}

// GetStatus 获取 pm2 状态
func (p *Pm2) GetStatus() (map[string]string, error) {
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

	// 获取 pm2 版本
	version := ""
	err, _ = p.RunShellWithOutput("pm2", "--version")
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

func (p *Pm2) GetInfo() core.SoftwareInfo {
	return p.info
}

func (p *Pm2) Start() error {
	p.Infof("pm2 无需启动服务")
	return nil
}

func (p *Pm2) Stop() error {
	p.Infof("pm2 无需停止服务")
	return nil
}
