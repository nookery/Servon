package softwares

// SoftwareManager 管理所有软件的安装、卸载等操作
type SoftwareManager struct {
	supportedSoftware []SoftwareInfo
}

// NewSoftwareManager 创建软件管理器实例
func NewSoftwareManager() *SoftwareManager {
	return &SoftwareManager{
		supportedSoftware: []SoftwareInfo{
			{Name: "nginx", Description: "高性能的 HTTP 和反向代理服务器"},
			{Name: "mongodb", Description: "流行的 NoSQL 数据库"},
			{Name: "redis", Description: "内存数据结构存储系统"},
			{Name: "mysql", Description: "流行的关系型数据库"},
			{Name: "docker", Description: "应用容器引擎"},
			{Name: "postgresql", Description: "开源对象关系数据库系统"},
			{Name: "nodejs", Description: "JavaScript 运行时环境"},
			{Name: "php", Description: "流行的服务端脚本语言"},
			{Name: "python3", Description: "Python 编程语言"},
			{Name: "golang", Description: "Go 编程语言"},
			{Name: "apache2", Description: "Apache HTTP 服务器"},
			{Name: "rabbitmq-server", Description: "开源消息队列系统"},
			{Name: "elasticsearch", Description: "分布式搜索和分析引擎"},
			{Name: "jenkins", Description: "开源持续集成工具"},
			{Name: "prometheus", Description: "系统监控和告警工具"},
		},
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
