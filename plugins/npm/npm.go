package npm

import (
	"fmt"
	"os/exec"
	"servon/core"
	"strings"
)

func Setup(app *core.App) {
	npm := NewNpm(app)
	app.RegisterSoftware("npm", npm)
}

// Npm 实现 Software 接口
type Npm struct {
	info core.SoftwareInfo
	*core.App
}

func NewNpm(app *core.App) core.SuperSoft {
	return &Npm{
		App: app,
		info: core.SoftwareInfo{
			Name:        "npm",
			Description: "Node.js 默认的包管理器",
		},
	}
}

func (n *Npm) Install() error {
	n.Infof("正在检查 npm...")

	// 检查 nodejs 是否已安装
	nodeCmd := exec.Command("node", "--version")
	if err := nodeCmd.Run(); err != nil {
		errMsg := "请先安装 NodeJS"
		return n.LogAndReturnErrorf("%s", errMsg)
	}

	// 检查 npm 是否已安装
	err, _ := n.RunShellWithOutput("npm", "--version")
	if err != nil {
		n.Infof("npm 未安装，正在通过 apt 安装...")

		// 使用 apt 安装 npm
		err, _ := n.RunShell("apt", "install", "-y", "npm")
		if err != nil {
			errMsg := "npm 安装失败"
			return n.LogAndReturnErrorf("%s: %v", errMsg, err)
		}
	}

	n.Success("npm 已安装")
	return nil
}

// Uninstall 卸载 npm
func (n *Npm) Uninstall() error {
	return fmt.Errorf("npm 是 NodeJS 的一部分，无法单独卸载")
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

func (n *Npm) GetInfo() core.SoftwareInfo {
	return n.info
}

func (n *Npm) Start() error {
	return fmt.Errorf("npm 是包管理工具，无需启动服务")
}

func (n *Npm) Stop() error {
	return fmt.Errorf("npm 是包管理工具，无需停止服务")
}
