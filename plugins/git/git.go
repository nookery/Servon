package git

import (
	"fmt"
	"os/exec"
	"servon/core"
	"strings"
)

func Setup(app *core.App) {
	app.RegisterSoftware("git", NewGit(app))
}

// Git 实现 Software 接口
type Git struct {
	info core.SoftwareInfo
	*core.App
}

func NewGit(app *core.App) core.SuperSoft {
	return &Git{
		App: app,
		info: core.SoftwareInfo{
			Name:        "git",
			Description: "分布式版本控制系统",
		},
	}
}

// Install 安装 Git
func (g *Git) Install() error {
	g.Infof("正在安装 Git...")

	// 检查操作系统类型
	osType := g.GetOSType()
	g.Infof("检测到操作系统: %s", osType)

	switch osType {
	case core.Ubuntu, core.Debian:
		// 更新软件包索引
		if err := g.AptUpdate(); err != nil {
			errMsg := fmt.Sprintf("更新软件包索引失败: %v", err)
			g.Errorf("%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

		// 安装 Git
		if err := g.AptInstall("git"); err != nil {
			errMsg := fmt.Sprintf("安装 Git 失败: %v", err)
			g.Errorf("%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

	case core.CentOS, core.RedHat:
		// 使用 yum 安装
		cmd := exec.Command("yum", "install", "-y", "git")
		err, _ := g.RunShell(cmd.String())
		if err != nil {
			errMsg := fmt.Sprintf("安装 Git 失败: %v", err)
			g.Errorf("%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

	default:
		errMsg := fmt.Sprintf("不支持的操作系统: %s", osType)
		g.Errorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	return nil
}

// Uninstall 卸载 Git
func (g *Git) Uninstall() error {
	g.Infof("正在卸载 Git...")

	osType := g.GetOSType()
	switch osType {
	case core.Ubuntu, core.Debian:
		if err := g.AptRemove("git"); err != nil {
			errMsg := fmt.Sprintf("卸载 Git 失败: %v", err)
			return g.LogAndReturnErrorf("%s", errMsg)
		}

	case core.CentOS, core.RedHat:
		cmd := exec.Command("yum", "remove", "-y", "git")
		err, _ := g.RunShell(cmd.String())
		if err != nil {
			errMsg := fmt.Sprintf("卸载 Git 失败: %v", err)
			return g.LogAndReturnErrorf("%s", errMsg)
		}

	default:
		errMsg := fmt.Sprintf("不支持的操作系统: %s", osType)
		return g.LogAndReturnErrorf("%s", errMsg)
	}

	g.Success("Git 卸载完成")
	return nil
}

func (g *Git) GetStatus() (map[string]string, error) {
	// 检查是否安装了 git
	cmd := exec.Command("git", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	// 解析版本信息
	version := strings.TrimSpace(string(output))
	return map[string]string{
		"status":  "installed",
		"version": version,
	}, nil
}

func (g *Git) GetInfo() core.SoftwareInfo {
	return g.info
}

func (g *Git) Start() error {
	g.Infof("Git 是版本控制工具，无需启动服务")
	return nil
}

func (g *Git) Stop() error {
	g.Infof("Git 是版本控制工具，无需停止服务")
	return nil
}
