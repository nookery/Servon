package contract

// Software 定义基础软件操作接口
type Software interface {
	// Install 安装软件，如果提供了日志通道则输出日志
	Install() error
	// Uninstall 卸载软件，如果提供了日志通道则输出日志
	Uninstall() error
	// GetStatus 获取软件状态
	GetStatus() (map[string]string, error)
	// Stop 停止软件服务
	Stop() error
	// Start 启动软件服务，如果提供了日志通道则输出日志
	Start() error
	// GetInfo 获取软件信息
	GetInfo() SoftwareInfo
}

// Gateway 定义网关特有的操作接口
type Gateway interface {
	GetConfig() (map[string]interface{}, error)
	SetConfig(config map[string]interface{}) error
	GetProjects() ([]Project, error)
	AddProject(project Project) error
	RemoveProject(projectName string) error
	ReloadConfig() error
}

// SuperSoft 组合基础软件和网关功能
type SuperSoft interface {
	Software
}

// SuperGateway 定义网关特有的操作接口
type SuperGateway interface {
	Software
	Gateway
}

// SoftwareInfo 软件基本信息
type SoftwareInfo struct {
	Name            string
	Description     string
	IsProxySoftware bool
	IsGateway       bool
}

// Project 网关项目配置
type Project struct {
	Name        string                 `json:"name"`
	Domain      string                 `json:"domain"`
	UpstreamURL string                 `json:"upstream_url"`
	Enabled     bool                   `json:"enabled"`
	Config      map[string]interface{} `json:"config"`
}
