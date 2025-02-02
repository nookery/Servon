package softwares

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
func (m *SoftwareManager) InstallSoftware(name string) (chan string, error) {
	sw, err := NewSoftware(name)
	if err != nil {
		return nil, err
	}
	return sw.Install()
}

// UninstallSoftware 卸载指定的软件
func (m *SoftwareManager) UninstallSoftware(name string) (chan string, error) {
	sw, err := NewSoftware(name)
	if err != nil {
		return nil, err
	}
	return sw.Uninstall()
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
