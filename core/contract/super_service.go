package contract

// SuperService 定义服务特有的操作接口
type SuperService interface {
	Software
	Service
}

// Service 定义服务特有的操作接口
type Service interface {
	// 启动软件本身
	Start() error

	// 启动某个服务
	StartService(serviceName string) error

	// 停止软件本身
	Stop() error

	// 停止某个服务
	StopService(serviceName string) error

	// 提供一个后台运行的命令，增加一个后台服务
	AddBackgroundService(serviceName string, command string, args []string, env []string) (string, error)

	// 停止后台运行的服务
	StopBackgroundService(serviceName string) error

	// 重启软件本身
	Restart() error

	// 重启某个服务
	RestartService(serviceName string) error

	// 获取服务日志
	GetLogs(serviceName string, lines int) (string, error)

	// 检查服务是否正在运行
	IsRunning(serviceName string) (bool, error)

	// 获取服务配置
	GetServiceConfig(serviceName string) (map[string]interface{}, error)

	// 更新服务配置
	UpdateServiceConfig(serviceName string, config map[string]interface{}) error

	// 获取服务资源使用情况
	GetResourceUsage(serviceName string) (map[string]interface{}, error)

	// 获取服务状态详情
	GetServiceDetails(serviceName string) (map[string]interface{}, error)
}
