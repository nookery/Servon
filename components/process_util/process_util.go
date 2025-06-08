// Package process_util 提供进程管理功能
package process_util

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

var (
	DefaultProcessUtil = NewProcessUtil()
)

type ProcessUtil struct{}

func NewProcessUtil() *ProcessUtil {
	return &ProcessUtil{}
}

// AutoStopPortProcess 自动停止占用指定端口的进程
func (p *ProcessUtil) AutoStopPortProcess(port int) error {
	// 使用 gopsutil 获取所有网络连接
	connections, err := net.Connections("tcp")
	if err != nil {
		return fmt.Errorf("failed to get network connections: %v", err)
	}

	// 查找使用指定端口的进程
	for _, conn := range connections {
		if conn.Laddr.Port == uint32(port) {
			// 获取进程信息
			proc, err := process.NewProcess(conn.Pid)
			if err != nil {
				continue
			}

			// 获取进程名称
			name, err := proc.Name()
			if err != nil {
				name = "unknown"
			}

			fmt.Printf("发现进程 %s (PID: %d) 占用端口 %d，正在终止...\n", name, conn.Pid, port)

			// 尝试优雅地终止进程
			err = proc.Terminate()
			if err != nil {
				// 如果优雅终止失败，强制杀死进程
				err = proc.Kill()
				if err != nil {
					return fmt.Errorf("failed to kill process %d: %v", conn.Pid, err)
				}
				fmt.Printf("强制终止进程 %s (PID: %d)\n", name, conn.Pid)
			} else {
				fmt.Printf("成功终止进程 %s (PID: %d)\n", name, conn.Pid)
			}
		}
	}

	return nil
}

// GetProcessByPort 获取占用指定端口的进程信息
func (p *ProcessUtil) GetProcessByPort(port int) ([]*process.Process, error) {
	var processes []*process.Process

	// 获取所有网络连接
	connections, err := net.Connections("tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get network connections: %v", err)
	}

	// 查找使用指定端口的进程
	for _, conn := range connections {
		if conn.Laddr.Port == uint32(port) {
			proc, err := process.NewProcess(conn.Pid)
			if err != nil {
				continue
			}
			processes = append(processes, proc)
		}
	}

	return processes, nil
}