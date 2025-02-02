package softwares

import "fmt"

// Software 定义软件操作的门面接口
type Software interface {
	// Install 安装软件并返回日志输出通道
	Install() (chan string, error)
	// Uninstall 卸载软件并返回日志输出通道
	Uninstall() (chan string, error)
	// GetStatus 获取软件状态
	GetStatus() (map[string]string, error)
	// Stop 停止软件服务
	Stop() error
	// GetInfo 获取软件信息
	GetInfo() SoftwareInfo
}

// SoftwareInfo 软件基本信息
type SoftwareInfo struct {
	Name        string
	Description string
}

// NewSoftware 创建指定软件的门面实例
func NewSoftware(name string) (Software, error) {
	switch name {
	case "nginx":
		return NewNginx(), nil
	case "mongodb":
		return NewMongoDB(), nil
	case "redis":
		return NewRedis(), nil
	case "mysql":
		return NewMySQL(), nil
	case "docker":
		return NewDocker(), nil
	// ... 其他软件
	default:
		return nil, fmt.Errorf("不支持的软件: %s", name)
	}
}
