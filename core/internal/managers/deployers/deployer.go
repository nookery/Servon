package deployers

import "servon/core/internal/utils"

// Deployer 定义了部署器的接口
type Deployer interface {
	// Deploy 部署项目
	Deploy(projectName string, workDir string, targetDir string, logger *utils.LogUtil) error

	// GetName 获取部署器名称
	GetName() string
}
