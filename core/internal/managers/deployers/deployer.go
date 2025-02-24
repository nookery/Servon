package deployers

import "servon/core/internal/utils"

// Deployer 定义了部署器的接口
type Deployer interface {
	// CanHandle 判断是否可以处理该项目
	CanHandle(workDir string) bool
	// Build 构建项目
	Build(workDir string, logger *utils.LogUtil) error
	// Deploy 部署项目
	Deploy(workDir string, targetDir string, logger *utils.LogUtil) error

	// GetName 获取部署器名称
	GetName() string
}
