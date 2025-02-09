package git

import (
	"fmt"
	"os/exec"
	"servon/core"
	"servon/core/contract"
	"servon/core/model"
	"servon/core/system"
	"servon/core/utils"
	"servon/core/utils/logger"
	"strings"
)

func Setup(core *core.Core) {
	core.RegisterSoftware("git", NewGit())
}

// Git 实现 Software 接口
type Git struct {
	info contract.SoftwareInfo
}

func NewGit() contract.SuperSoft {
	return &Git{
		info: contract.SoftwareInfo{
			Name:        "git",
			Description: "分布式版本控制系统",
		},
	}
}

// ... existing code from git.go ...
func (g *Git) Install(logChan chan<- string) error {
	logger.InfoChan(logChan, "正在安装 Git...")

	// 检查操作系统类型
	osType := utils.GetOSType()
	logger.InfoChan(logChan, "检测到操作系统: %s", osType)

	switch osType {
	case model.Ubuntu, model.Debian:
		apt := system.NewApt()

		// 更新软件包索引
		if err := apt.Update(); err != nil {
			errMsg := fmt.Sprintf("更新软件包索引失败: %v", err)
			logger.ErrorChan(logChan, "%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

		// 安装 Git
		if err := apt.Install("git"); err != nil {
			errMsg := fmt.Sprintf("安装 Git 失败: %v", err)
			logger.ErrorChan(logChan, "%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

	case model.CentOS, model.RedHat:
		// 使用 yum 安装
		cmd := exec.Command("yum", "install", "-y", "git")
		if err := logger.StreamCommand(cmd); err != nil {
			errMsg := fmt.Sprintf("安装 Git 失败: %v", err)
			logger.ErrorChan(logChan, "%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

	default:
		errMsg := fmt.Sprintf("不支持的操作系统: %s", osType)
		logger.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	logger.InfoChan(logChan, "Git 安装完成")
	return nil
}

func (g *Git) Uninstall(logChan chan<- string) error {
	logger.InfoChan(logChan, "正在卸载 Git...")

	osType := utils.GetOSType()
	switch osType {
	case model.Ubuntu, model.Debian:
		apt := system.NewApt()
		if err := apt.Remove("git"); err != nil {
			errMsg := fmt.Sprintf("卸载 Git 失败: %v", err)
			logger.ErrorChan(logChan, "%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

	case model.CentOS, model.RedHat:
		cmd := exec.Command("yum", "remove", "-y", "git")
		if err := logger.StreamCommand(cmd); err != nil {
			errMsg := fmt.Sprintf("卸载 Git 失败: %v", err)
			logger.ErrorChan(logChan, "%s", errMsg)
			return fmt.Errorf("%s", errMsg)
		}

	default:
		errMsg := fmt.Sprintf("不支持的操作系统: %s", osType)
		logger.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	logger.InfoChan(logChan, "Git 卸载完成")
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

func (g *Git) GetInfo() contract.SoftwareInfo {
	return g.info
}

func (g *Git) Start(logChan chan<- string) error {
	logger.InfoChan(logChan, "Git 是版本控制工具，无需启动服务")
	return nil
}

func (g *Git) Stop() error {
	logger.Info("Git 是版本控制工具，无需停止服务")
	return nil
}
