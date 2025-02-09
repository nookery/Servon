package provider

import (
	"github.com/spf13/cobra"
	"servon/core/model"
	"servon/core/utils"
	"servon/core/utils/logger"
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

func (s *SystemProvider) GetOSType() utils.OSType {
	return utils.GetOSType()
}

func (s *SystemProvider) CanUseApt() bool {
	osType := s.GetOSType()
	return osType == model.Ubuntu || osType == model.Debian
}
