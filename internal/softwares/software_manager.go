package softwares

import "fmt"

// SoftwareManager 管理所有软件的安装、卸载等操作
type SoftwareManager struct {
	supportedSoftware []SoftwareInfo
}

// NewSoftwareManager 创建软件管理器实例
func NewSoftwareManager() *SoftwareManager {
	return &SoftwareManager{
		supportedSoftware: GetSupportedSoftware(),
	}
}

// GetSupportedSoftware 返回支持的软件列表
func (m *SoftwareManager) GetSupportedSoftware() []SoftwareInfo {
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

// InstallSoftware 安装指定的软件
func (m *SoftwareManager) InstallSoftware(name string, msgChan chan<- string) error {
	fmt.Println("安装软件, 是否提供管道", msgChan != nil)
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
func (m *SoftwareManager) StartSoftware(name string) (chan string, error) {
	sw, err := NewSoftware(name)
	if err != nil {
		return nil, err
	}

	// 创建一个新的 channel 来包装原始的 channel
	outputChan := make(chan string, 100)
	err = sw.Start(outputChan)
	if err != nil {
		return nil, err
	}

	return outputChan, nil
}
