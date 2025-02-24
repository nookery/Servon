package nodejs

import (
	"fmt"
	"os/exec"
	"servon/core"
	"strings"
)

type NodeJSPlugin struct {
	info core.SoftwareInfo
	*core.App
}

func Setup(app *core.App) {
	nodejs := NewNodeJSPlugin(app)
	app.RegisterSoftware("nodejs", nodejs)
}

func NewNodeJSPlugin(app *core.App) core.SuperSoft {
	return &NodeJSPlugin{
		App: app,
		info: core.SoftwareInfo{
			Name:        "nodejs",
			Description: "JavaScript 运行时环境",
		},
	}
}

func (n *NodeJSPlugin) Install() error {
	osType := n.GetOSType()

	switch osType {
	case core.Ubuntu, core.Debian:
		n.Infof("使用 APT 包管理器安装...")
		n.Infof("添加 NodeJS 官方源...")

		// 下载并安装 NodeSource 设置脚本
		err, _ := n.RunShell("sh", "-c", "curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -")
		if err != nil {
			return err
		}

		// 安装 NodeJS
		if err := n.AptInstall("nodejs"); err != nil {
			return err
		}

	case core.CentOS, core.RedHat:
		errMsg := "暂不支持在 RHEL 系统上安装 NodeJS"
		return n.LogAndReturnErrorf("%s", errMsg)

	default:
		errMsg := fmt.Sprintf("不支持的操作系统: %s", osType)
		return n.LogAndReturnErrorf("%s", errMsg)
	}

	// 验证安装结果
	if !n.IsInstalled("nodejs") {
		errMsg := "NodeJS 安装验证失败，未检测到已安装的包"
		return n.LogAndReturnErrorf("%s", errMsg)
	}

	n.Success("NodeJS 安装完成")
	return nil
}

// Uninstall 卸载 NodeJS
func (n *NodeJSPlugin) Uninstall() error {
	// 卸载软件包及其依赖
	if err := n.AptRemove("nodejs"); err != nil {
		return err
	}

	// 清理配置文件
	if err := n.AptPurge("nodejs"); err != nil {
		return err
	}

	// 删除 NodeSource 源文件
	err, _ := n.RunShell("sudo", "rm", "/etc/apt/sources.list.d/nodesource.list")
	if err != nil {
		return fmt.Errorf("删除源文件失败:\n%s", err)
	}

	// 清理自动安装的依赖
	err, _ = n.RunShell("sudo", "apt-get", "autoremove", "-y")
	if err != nil {
		return fmt.Errorf("清理依赖失败:\n%s", err)
	}

	n.Success("NodeJS 卸载完成")
	return nil
}

func (n *NodeJSPlugin) GetStatus() (map[string]string, error) {
	if !n.IsInstalled("nodejs") {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	// 获取版本
	version := ""
	verCmd := exec.Command("node", "--version")
	if verOutput, err := verCmd.CombinedOutput(); err == nil {
		version = strings.TrimSpace(string(verOutput))
	}

	return map[string]string{
		"status":  "installed",
		"version": version,
	}, nil
}

func (n *NodeJSPlugin) GetInfo() core.SoftwareInfo {
	return n.info
}

func (n *NodeJSPlugin) Start() error {
	n.Infof("NodeJS 是运行时环境，无需启动服务")
	return nil
}

func (n *NodeJSPlugin) Stop() error {
	n.Infof("NodeJS 是运行时环境，无需停止服务")
	return nil
}
