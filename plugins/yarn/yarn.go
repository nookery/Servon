package yarn

import (
	"fmt"
	"servon/core"
	"servon/core/contract"
	"strings"
)

// Yarn 实现 Software 接口
type Yarn struct {
	info contract.SoftwareInfo
	*core.Core
}

func Setup(core *core.Core) {
	core.RegisterSoftware("yarn", NewYarn(core))
}

func NewYarn(core *core.Core) contract.SuperSoft {
	return &Yarn{
		Core: core,
		info: contract.SoftwareInfo{
			Name:        "yarn",
			Description: "快速、可靠、安全的依赖管理工具",
		},
	}
}

func (y *Yarn) Install(logChan chan<- string) error {
	y.PrintInfo("正在安装 Yarn...")

	// 检查 nodejs 是否已安装
	if _, err := y.RunShellWithOutput("node", "--version"); err != nil {
		errMsg := "请先安装 NodeJS"
		y.PrintErrorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 使用 StreamCommand 来执行安装并输出详细日志
	if err := y.RunShell("npm", "install", "-g", "yarn"); err != nil {
		errMsg := fmt.Sprintf("安装 Yarn 失败: %v", err)
		y.PrintErrorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	y.PrintSuccess("Yarn 安装完成")
	return nil
}

func (y *Yarn) Uninstall(logChan chan<- string) error {
	y.PrintInfo("正在卸载 Yarn...")

	if err := y.RunShell("npm", "uninstall", "-g", "yarn"); err != nil {
		errMsg := fmt.Sprintf("卸载 Yarn 失败: %v", err)
		y.PrintErrorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	y.PrintSuccess("Yarn 卸载完成")
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

func (y *Yarn) GetInfo() contract.SoftwareInfo {
	return y.info
}

func (y *Yarn) Start(logChan chan<- string) error {
	y.PrintInfo("Yarn 是包管理工具，无需启动服务")
	return nil
}

func (y *Yarn) Stop() error {
	y.PrintInfo("Yarn 是包管理工具，无需停止服务")
	return nil
}
