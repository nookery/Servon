package nodejs

import (
	"fmt"
	"os/exec"
	"servon/core"
	"servon/core/contract"
	"strings"
)

type NodeJSPlugin struct {
	info contract.SoftwareInfo
	*core.Core
}

func Setup(core *core.Core) {
	nodejs := NewNodeJSPlugin(core)
	core.RegisterSoftware("nodejs", nodejs)
}

func NewNodeJSPlugin(core *core.Core) contract.SuperSoft {
	return &NodeJSPlugin{
		Core: core,
		info: contract.SoftwareInfo{
			Name:        "nodejs",
			Description: "JavaScript 运行时环境",
		},
	}
}

func (n *NodeJSPlugin) Install(logChan chan<- string) error {
	outputChan := make(chan string, 100)
	osType := n.GetOSType()

	switch osType {
	case core.Ubuntu, core.Debian:
		n.PrintInfo("使用 APT 包管理器安装...")
		n.PrintInfo("添加 NodeJS 官方源...")

		// 下载并安装 NodeSource 设置脚本
		if err := n.RunShell("sh", "-c", "curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -"); err != nil {
			return err
		}

		// 安装 NodeJS
		if err := n.AptInstall("nodejs"); err != nil {
			return err
		}

	case core.CentOS, core.RedHat:
		errMsg := "暂不支持在 RHEL 系统上安装 NodeJS"
		n.ErrorChan(outputChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)

	default:
		errMsg := fmt.Sprintf("不支持的操作系统: %s", osType)
		n.ErrorChan(outputChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 验证安装结果
	if !n.IsInstalled("nodejs") {
		errMsg := "NodeJS 安装验证失败，未检测到已安装的包"
		n.ErrorChan(outputChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	n.PrintSuccess("NodeJS 安装完成")
	return nil
}

func (n *NodeJSPlugin) Uninstall(logChan chan<- string) error {
	outputChan := make(chan string, 100)

	go func() {
		defer close(outputChan)

		// 卸载软件包及其依赖
		if err := n.AptRemove("nodejs"); err != nil {
			return
		}

		// 清理配置文件
		if err := n.AptPurge("nodejs"); err != nil {
			return
		}

		// 删除 NodeSource 源文件
		rmSourceCmd := exec.Command("sudo", "rm", "/etc/apt/sources.list.d/nodesource.list")
		if output, err := rmSourceCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("删除源文件失败:\n%s", string(output))
			return
		}

		// 清理自动安装的依赖
		cleanCmd := exec.Command("sudo", "apt-get", "autoremove", "-y")
		if output, err := cleanCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("清理依赖失败:\n%s", string(output))
			return
		}

		outputChan <- "NodeJS 卸载完成"
	}()

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

func (n *NodeJSPlugin) GetInfo() contract.SoftwareInfo {
	return n.info
}

func (n *NodeJSPlugin) Start(logChan chan<- string) error {
	n.Info("NodeJS 是运行时环境，无需启动服务")
	return nil
}

func (n *NodeJSPlugin) Stop() error {
	n.Info("NodeJS 是运行时环境，无需停止服务")
	return nil
}
