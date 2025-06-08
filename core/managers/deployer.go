package managers

import "servon/components/log_util"

// Deployer 定义了部署器的接口
type Deployer interface {
	// Deploy 部署项目
	Deploy(projectName string, workDir string, targetDir string, logger *log_util.LogUtil) error

	// GetName 获取部署器名称
	GetName() string
}
