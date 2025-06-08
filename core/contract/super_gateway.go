package contract

// SuperGateway 定义网关特有的操作接口
type SuperGateway interface {
	Software
	Gateway
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

// Project 网关项目配置
type Project struct {
	Name        string                 `json:"name"`
	Domain      string                 `json:"domain"`
	UpstreamURL string                 `json:"upstream_url"`
	Enabled     bool                   `json:"enabled"`
	Config      map[string]interface{} `json:"config"`
}
