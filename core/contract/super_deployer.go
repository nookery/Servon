package contract

import logger1 "servon/components/logger"

// SuperDeployer 定义了部署器的接口
type SuperDeployer interface {
	// Deploy 部署项目
	Deploy(projectName string, workDir string, targetDir string, logger *logger1.LogUtil) error

	// GetName 获取部署器名称
	GetName() string
}
