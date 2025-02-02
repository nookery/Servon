package softwares

import "fmt"

// SoftwareRegistry 存储所有已注册的软件
var registry = map[string]func() Software{
	"nginx":   func() Software { return NewNginx() },
	"mongodb": func() Software { return NewMongoDB() },
	"redis":   func() Software { return NewRedis() },
	"mysql":   func() Software { return NewMySQL() },
	"docker":  func() Software { return NewDocker() },
}

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
	if factory, ok := registry[name]; ok {
		return factory(), nil
	}
	return nil, fmt.Errorf("不支持的软件: %s", name)
}

// GetSupportedSoftware 获取所有支持的软件
func GetSupportedSoftware() []SoftwareInfo {
	var list []SoftwareInfo
	for _, factory := range registry {
		sw := factory()
		list = append(list, sw.GetInfo())
	}
	return list
}
