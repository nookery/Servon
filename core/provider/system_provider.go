package provider

import (
	"os/exec"
	"servon/core/model"
	"servon/core/utils/logger"
	"strings"

	"github.com/spf13/cobra"
)

// SystemProvider 系统管理器
type SystemProvider struct {
	RootCmd *cobra.Command
}

func NewSystemProvider() SystemProvider {
	return SystemProvider{
		RootCmd: &cobra.Command{},
	}
}

func (p *SystemProvider) InstallSoftware(name string, logChan chan<- string) error {
	// 检查操作系统类型
	osType := p.GetOSType()
	logger.InfoChan(logChan, "检测到操作系统: %s", osType)

	return nil
}

func (p *SystemProvider) UninstallSoftware(name string, logChan chan<- string) error {
	return nil
}

func (s *SystemProvider) GetOSType() OSType {
	// 尝试读取 /etc/os-release 文件
	cmd := exec.Command("cat", "/etc/os-release")
	output, err := cmd.Output()
	if err != nil {
		return model.Unknown
	}

	osInfo := strings.ToLower(string(output))

	switch {
	case strings.Contains(osInfo, "ubuntu"):
		return model.Ubuntu
	case strings.Contains(osInfo, "debian"):
		return model.Debian
	case strings.Contains(osInfo, "centos"):
		return model.CentOS
	case strings.Contains(osInfo, "redhat"):
		return model.RedHat
	default:
		return model.Unknown
	}
}

func (s *SystemProvider) CanUseApt() bool {
	osType := s.GetOSType()
	return osType == model.Ubuntu || osType == model.Debian
}
