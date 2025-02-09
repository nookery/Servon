package provider

import (
	"fmt"
	"servon/core/contract"
)

type SuperSoft = contract.SuperSoft

// SoftwareProvider 软件提供者
type SoftwareProvider struct {
	softwares map[string]SuperSoft
}

// NewSoftwareProvider 创建软件提供者
func NewSoftwareProvider() SoftwareProvider {
	return SoftwareProvider{
		softwares: make(map[string]SuperSoft),
	}
}

// Register 注册软件
func (p *SoftwareProvider) Register(name string, software SuperSoft) error {
	if _, exists := p.softwares[name]; exists {
		return fmt.Errorf("software %s already registered", name)
	}
	p.softwares[name] = software
	return nil
}

// Install 安装软件, 如果提供了日志通道则输出日志
func (p *SoftwareProvider) Install(name string, logChan chan<- string) error {
	software, ok := p.softwares[name]
	if !ok {
		return fmt.Errorf("software %s not found", name)
	}
	return software.Install(logChan)
}

// Uninstall 卸载软件, 如果提供了日志通道则输出日志
func (p *SoftwareProvider) Uninstall(name string, logChan chan<- string) error {
	software, ok := p.softwares[name]
	if !ok {
		return fmt.Errorf("software %s not found", name)
	}
	return software.Uninstall(logChan)
}

// Start 启动软件, 如果提供了日志通道则输出日志
func (p *SoftwareProvider) Start(name string, logChan chan<- string) error {
	software, ok := p.softwares[name]
	if !ok {
		return fmt.Errorf("software %s not found", name)
	}
	return software.Start(logChan)
}

// Stop 停止软件, 如果提供了日志通道则输出日志
func (p *SoftwareProvider) Stop(name string) error {
	software, ok := p.softwares[name]
	if !ok {
		return fmt.Errorf("software %s not found", name)
	}
	return software.Stop()
}

// GetStatus 获取软件状态
func (p *SoftwareProvider) GetStatus(name string) (map[string]string, error) {
	software, ok := p.softwares[name]
	if !ok {
		return nil, fmt.Errorf("software %s not found", name)
	}
	return software.GetStatus()
}

// GetAllSoftware 获取所有软件
func (p *SoftwareProvider) GetAllSoftware() []string {
	softwareNames := make([]string, 0, len(p.softwares))
	for name := range p.softwares {
		softwareNames = append(softwareNames, name)
	}
	return softwareNames
}
