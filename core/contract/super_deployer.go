package contract

// SuperDeployer 定义了部署器的接口
type SuperDeployer interface {
	// Deploy 部署项目
	Deploy(projectName string, workDir string, targetDir string) error

	// GetName 获取部署器名称
	GetName() string
}
