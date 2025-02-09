package yarn

import (
	"fmt"
	"os/exec"
	"servon/cmd/software"
	"servon/cmd/utils/logger"
	"servon/core/contract"
	"strings"
)

// YarnPlugin 实现 Plugin 接口
type YarnPlugin struct{}

func (p *YarnPlugin) Init() error {
	return nil
}

func (p *YarnPlugin) Name() string {
	return "yarn"
}

func (p *YarnPlugin) Register() {
	software.RegisterSoftware("yarn", func() contract.SuperSoftware {
		return NewYarn()
	})
}

// Yarn 实现 Software 接口
type Yarn struct {
	info contract.SoftwareInfo
}

func NewYarn() contract.SuperSoftware {
	return &Yarn{
		info: contract.SoftwareInfo{
			Name:        "yarn",
			Description: "快速、可靠、安全的依赖管理工具",
		},
	}
}

func (y *Yarn) Install(logChan chan<- string) error {
	logger.InfoChan(logChan, "正在安装 Yarn...")

	// 检查 nodejs 是否已安装
	nodeCmd := exec.Command("node", "--version")
	if err := nodeCmd.Run(); err != nil {
		errMsg := "请先安装 NodeJS"
		logger.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 使用 StreamCommand 来执行安装并输出详细日志
	cmd := exec.Command("npm", "install", "-g", "yarn")
	if err := logger.StreamCommand(cmd); err != nil {
		errMsg := fmt.Sprintf("安装 Yarn 失败: %v", err)
		logger.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	logger.InfoChan(logChan, "Yarn 安装完成")
	return nil
}

func (y *Yarn) Uninstall(logChan chan<- string) error {
	logger.InfoChan(logChan, "正在卸载 Yarn...")

	cmd := exec.Command("npm", "uninstall", "-g", "yarn")
	if err := logger.StreamCommand(cmd); err != nil {
		errMsg := fmt.Sprintf("卸载 Yarn 失败: %v", err)
		logger.ErrorChan(logChan, "%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	logger.InfoChan(logChan, "Yarn 卸载完成")
	return nil
}

func (y *Yarn) GetStatus() (map[string]string, error) {
	// 检查 nodejs 是否已安装
	nodeCmd := exec.Command("node", "--version")
	if err := nodeCmd.Run(); err != nil {
		return map[string]string{
			"status":  "nodejs_not_installed",
			"version": "",
		}, nil
	}

	// 获取 yarn 版本
	version := ""
	verCmd := exec.Command("yarn", "--version")
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

func (y *Yarn) GetInfo() contract.SoftwareInfo {
	return y.info
}

func (y *Yarn) Start(logChan chan<- string) error {
	logger.InfoChan(logChan, "Yarn 是包管理工具，无需启动服务")
	return nil
}

func (y *Yarn) Stop() error {
	logger.Info("Yarn 是包管理工具，无需停止服务")
	return nil
}

func init() {
	// 在包被导入时自动注册插件
	if err := contract.RegisterPlugin(&YarnPlugin{}); err != nil {
		fmt.Printf("Failed to register Yarn plugin: %v\n", err)
	}
}
