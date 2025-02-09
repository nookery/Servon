package software

import (
	"fmt"
	"sync"

	"servon/core/contract"
)

var (
	// registry 存储所有已注册的软件
	registry = make(map[string]func() contract.SuperSoftware)
	regMutex sync.RWMutex
)

// RegisterSoftware 注册一个新的软件到注册表中
func RegisterSoftware(name string, factory func() contract.SuperSoftware) {
	regMutex.Lock()
	defer regMutex.Unlock()
	registry[name] = factory
}

// NewSoftware 创建一个新的软件实例
func NewSoftware(name string) (contract.SuperSoftware, error) {
	regMutex.RLock()
	factory, exists := registry[name]
	regMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("software %s is not supported", name)
	}
	return factory(), nil
}

// GetSupportedSoftware 返回支持的软件列表
func GetSupportedSoftware() []contract.SoftwareInfo {
	regMutex.RLock()
	defer regMutex.RUnlock()

	supportedSoftware := make([]contract.SoftwareInfo, 0, len(registry))
	for name, factory := range registry {
		sw := factory()
		info := contract.SoftwareInfo{
			Name:        name,
			Description: sw.GetInfo().Description,
		}
		supportedSoftware = append(supportedSoftware, info)
	}
	return supportedSoftware
}

// SoftwareManager 管理所有软件的安装、卸载等操作
type SoftwareManager struct {
	supportedSoftware []contract.SoftwareInfo
}

// NewSoftwareManager 创建软件管理器实例
func NewSoftwareManager() *SoftwareManager {
	return &SoftwareManager{
		supportedSoftware: GetSupportedSoftware(),
	}
}

// GetSupportedSoftware 返回支持的软件列表
func (m *SoftwareManager) GetSupportedSoftware() []contract.SoftwareInfo {
	return m.supportedSoftware
}

// GetSoftwareNames 返回支持的软件名称列表
func (m *SoftwareManager) GetSoftwareNames() []string {
	names := make([]string, len(m.supportedSoftware))
	for i, sw := range m.supportedSoftware {
		names[i] = sw.Name
	}
	return names
}

// IsSupportedSoftware checks if the given software is supported
func (sm *SoftwareManager) IsSupportedSoftware(name string) bool {
	for _, sw := range sm.GetSoftwareNames() {
		if sw == name {
			return true
		}
	}
	return false
}

// InstallSoftware 安装指定的软件
func (m *SoftwareManager) InstallSoftware(name string, msgChan chan<- string) error {
	sw, err := NewSoftware(name)
	if err != nil {
		return err
	}

	err = sw.Install(msgChan)
	if err != nil {
		return err
	}

	return nil
}

// UninstallSoftware 卸载指定的软件
func (m *SoftwareManager) UninstallSoftware(name string, msgChan chan<- string) error {
	sw, err := NewSoftware(name)
	if err != nil {
		return err
	}
	return sw.Uninstall(msgChan)
}

// GetSoftwareStatus 获取软件状态
func (m *SoftwareManager) GetSoftwareStatus(name string) (map[string]string, error) {
	sw, err := NewSoftware(name)
	if err != nil {
		return nil, err
	}
	return sw.GetStatus()
}

// StopSoftware 停止软件服务
func (m *SoftwareManager) StopSoftware(name string) error {
	sw, err := NewSoftware(name)
	if err != nil {
		return err
	}
	return sw.Stop()
}

// StartSoftware 启动指定的软件
func (m *SoftwareManager) StartSoftware(name string, logChan chan<- string) error {
	sw, err := NewSoftware(name)
	if err != nil {
		return err
	}

	return sw.Start(logChan)
}
