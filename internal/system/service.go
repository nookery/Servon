package system

import "servon/internal/system/service"

// ServiceIsActive 检查服务是否处于活动状态
func ServiceIsActive(serviceName string) bool {
	return service.IsActive(serviceName)
}

// ServiceStop 停止指定服务
func ServiceStop(serviceName string) error {
	return service.Stop(serviceName)
}

// ServiceReload 重新加载指定服务
func ServiceReload(serviceName string) error {
	return service.Reload(serviceName)
}

// ServiceStart 启动指定服务
func ServiceStart(serviceName string) error {
	return service.Start(serviceName)
}
