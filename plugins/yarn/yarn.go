package yarn

import (
	"fmt"
	"servon/core"
	"strings"
)

// Yarn 实现 Software 接口
type Yarn struct {
	info core.SoftwareInfo
	*core.App
}

func Setup(app *core.App) {
	app.RegisterSoftware("yarn", NewYarn(app))
}

func NewYarn(app *core.App) core.SuperSoft {
	return &Yarn{
		App: app,
		info: core.SoftwareInfo{
			Name:        "yarn",
			Description: "快速、可靠、安全的依赖管理工具",
		},
	}
}

func (y *Yarn) Install() error {
	y.Infof("正在安装 Yarn...")

	// 检查 nodejs 是否已安装
	if _, err := y.RunShellWithOutput("node", "--version"); err != nil {
		errMsg := "请先安装 NodeJS"
		return y.LogAndReturnErrorf("%s", errMsg)
	}

	// 使用 StreamCommand 来执行安装并输出详细日志
	if err := y.RunShell("npm", "install", "-g", "yarn"); err != nil {
		errMsg := fmt.Sprintf("安装 Yarn 失败: %v", err)
		return y.LogAndReturnErrorf("%s", errMsg)
	}

	y.Success("Yarn 安装完成")
	return nil
}

// Uninstall 卸载 Yarn
func (y *Yarn) Uninstall() error {
	y.Infof("正在卸载 Yarn...")

	if err := y.RunShell("npm", "uninstall", "-g", "yarn"); err != nil {
		errMsg := fmt.Sprintf("卸载 Yarn 失败: %v", err)
		return y.LogAndReturnErrorf("%s", errMsg)
	}

	y.Success("Yarn 卸载完成")
	return nil
}

func (y *Yarn) GetStatus() (map[string]string, error) {
	// 检查 nodejs 是否已安装
	if _, err := y.RunShellWithOutput("node", "--version"); err != nil {
		return map[string]string{
			"status":  "nodejs_not_installed",
			"version": "",
		}, nil
	}

	// 获取 yarn 版本
	version := ""
	if output, err := y.RunShellWithOutput("yarn", "--version"); err == nil {
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

func (y *Yarn) GetInfo() core.SoftwareInfo {
	return y.info
}

func (y *Yarn) Start() error {
	y.Infof("Yarn 是包管理工具，无需启动服务")
	return nil
}

func (y *Yarn) Stop() error {
	y.Infof("Yarn 是包管理工具，无需停止服务")
	return nil
}
