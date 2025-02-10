package api

import (
	"fmt"
	"servon/core/libs"
)

type SystemApi struct{}

func NewSystemApi() SystemApi {
	return SystemApi{}
}

func (c *SystemApi) UninstallSoftwareUseSystem(name string) error {
	return fmt.Errorf("not implemented")
}

func (c *SystemApi) GetOSType() libs.OSType {
	return libs.GetOSType()
}

func (c *SystemApi) CanUseApt() bool {
	return false
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
//   - serviceFilePath: 服务文件的路径
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
func (c *SystemApi) RunBackgroundService(command string, args []string, logChan chan<- string) (string, error) {
	return libs.RunBackgroundService(command, args, logChan)
}

// StopBackgroundService 停止并移除后台运行的服务
func (c *SystemApi) StopBackgroundService(command string, logChan chan<- string) error {
	return libs.StopBackgroundService(command, logChan)
}

// IsInstalled 检查软件包是否已安装
func (c *SystemApi) IsInstalled(name string) bool {
	dpkg := libs.NewDpkg(nil)
	return dpkg.IsInstalled(name)
}

// ServiceIsActive 检查服务是否正在运行
func (c *SystemApi) ServiceIsActive(name string) bool {
	return libs.IsActive(name)
}

// AptUpdate 更新软件包索引
func (c *SystemApi) AptUpdate() error {
	return libs.AptUpdate()
}

// AptInstall 安装指定的软件包
func (c *SystemApi) AptInstall(packages ...string) error {
	return libs.AptInstall(packages...)
}

// AptRemove 移除指定的软件包
func (c *SystemApi) AptRemove(packages ...string) error {
	return libs.AptRemove(packages...)
}

// AptPurge 完全移除软件包及其配置文件
func (c *SystemApi) AptPurge(packages ...string) error {
	return libs.AptPurge(packages...)
}

// ServiceReload 重载服务
func (c *SystemApi) ServiceReload(serviceName string) error {
	return libs.Reload(serviceName)
}

// ServiceStart 启动服务
func (c *SystemApi) ServiceStart(serviceName string) error {
	return libs.Start(serviceName)
}

// ServiceStop 停止服务
func (c *SystemApi) ServiceStop(serviceName string) error {
	return libs.Stop(serviceName)
}
