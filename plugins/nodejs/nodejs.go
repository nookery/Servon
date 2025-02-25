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
		n.SoftwareLogger.Infof("使用 APT 包管理器安装...")
		n.SoftwareLogger.Infof("添加 NodeJS 官方源...")

		// 首先检查sudo是否可用
		hasSudo := n.checkSudoAvailable()
		n.SoftwareLogger.Infof("系统sudo可用性: %v", hasSudo)

		n.SoftwareLogger.Infof("开始下载 NodeSource 设置脚本...")
		var err error
		var output string

		if hasSudo {
			// 使用sudo方式
			err, output = n.RunShell("sh", "-c", "curl -m 60 -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -")
		} else {
			// 无sudo环境（可能在容器或已是root用户）
			err, output = n.RunShell("sh", "-c", "curl -m 60 -fsSL https://deb.nodesource.com/setup_lts.x | bash -")
		}

		n.SoftwareLogger.Infof("NodeSource 脚本执行输出: %s", output)
		if err != nil {
			return n.SoftwareLogger.LogAndReturnErrorf("下载或执行 NodeSource 脚本失败: %v", err)
		}

		// 更新软件包列表
		n.SoftwareLogger.Infof("更新软件包列表...")
		if hasSudo {
			err, output = n.RunShell("sudo", "apt-get", "update", "-y")
		} else {
			err, output = n.RunShell("apt-get", "update", "-y")
		}

		if err != nil {
			n.SoftwareLogger.Errorf("更新软件包列表失败: %v, 输出: %s", err, output)
			return err
		}

		// 安装 NodeJS (使用-y参数确保非交互式)
		n.SoftwareLogger.Infof("开始安装 NodeJS...")
		err = n.AptInstall("nodejs")

		if err != nil {
			return n.SoftwareLogger.LogAndReturnErrorf("安装 NodeJS 失败: %v", err)
		}

		n.SoftwareLogger.Infof("NodeJS 安装完成，正在验证...")

	case core.CentOS, core.RedHat:
		errMsg := "暂不支持在 RHEL 系统上安装 NodeJS"
		return n.SoftwareLogger.LogAndReturnErrorf("%s", errMsg)

	default:
		errMsg := fmt.Sprintf("不支持的操作系统: %s", osType)
		return n.SoftwareLogger.LogAndReturnErrorf("%s", errMsg)
	}

	// 验证安装结果
	if !n.IsInstalled("nodejs") {
		errMsg := "NodeJS 安装验证失败，未检测到已安装的包"
		return n.SoftwareLogger.LogAndReturnErrorf("%s", errMsg)
	}

	n.SoftwareLogger.Success("NodeJS 安装完成")
	return nil
}

// Uninstall 卸载 NodeJS
func (n *NodeJSPlugin) Uninstall() error {
	// 卸载软件包及其依赖
	if err := n.AptRemove("nodejs"); err != nil {
		return n.SoftwareLogger.LogAndReturnError(err.Error())
	}

	// 清理配置文件
	if err := n.AptPurge("nodejs"); err != nil {
		return n.SoftwareLogger.LogAndReturnError(err.Error())
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

	n.SoftwareLogger.Success("NodeJS 卸载完成")
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
	n.SoftwareLogger.Infof("NodeJS 是运行时环境，无需启动服务")
	return nil
}

func (n *NodeJSPlugin) Stop() error {
	n.SoftwareLogger.Infof("NodeJS 是运行时环境，无需停止服务")
	return nil
}

// checkSudoAvailable 检查系统中是否有sudo命令
func (n *NodeJSPlugin) checkSudoAvailable() bool {
	err, _ := n.RunShell("which", "sudo")
	return err == nil
}
