package api

import (
	"servon/core/model"
	"servon/core/provider"
)

type System struct {
	systemProvider provider.SystemProvider
}

func NewSystem() System {
	return System{
		systemProvider: provider.NewSystemProvider(),
	}
}

func (c *System) UninstallSoftwareUseSystem(name string) error {
	return c.systemProvider.UninstallSoftware(name, nil)
}

func (c *System) GetOSType() model.OSType {
	return c.systemProvider.GetOSType()
}

func (c *System) CanUseApt() bool {
	return c.systemProvider.CanUseApt()
}

// RunBackgroundService 使用 systemd 在后台运行指定的命令作为服务
//
// 参数:
//   - command: 要执行的命令名称（例如 "node", "python3"）
//   - args: 命令的参数列表（例如 []string{"app.js", "--port", "3000"}）
//   - logChan: 日志通道，用于接收服务的运行状态信息
//
// 返回值:
//   - error: 如果服务启动失败，返回错误信息；如果成功启动，返回 nil
//
// 示例:
//
//	// 运行 Node.js 应用作为服务
//	logChan := make(chan string)
//	go func() {
//	    for msg := range logChan {
//	        fmt.Println(msg)
//	    }
//	}()
//	err := core.RunBackgroundService("node", []string{"app.js"}, logChan)
//
// 注意事项:
//  1. 需要 root 权限才能创建和管理 systemd 服务
//  2. 服务会在系统启动时自动启动
//  3. 如果进程崩溃会自动重启
//  4. 日志会被写入到 /var/log/servon/{服务名}.log
//  5. 可以使用 systemctl status/start/stop/restart 命令管理服务
func (c *System) RunBackgroundService(command string, args []string, logChan chan<- string) error {
	return c.systemProvider.RunBackgroundService(command, args, logChan)
}

// StopBackgroundService 停止并移除后台运行的服务
func (c *System) StopBackgroundService(command string, logChan chan<- string) error {
	return c.systemProvider.StopBackgroundService(command, logChan)
}
