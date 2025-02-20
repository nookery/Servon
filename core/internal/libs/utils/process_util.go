package utils

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
			if conn.Pid == 0 {
				continue
			}

			// 获取进程实例
			proc, err := process.NewProcess(conn.Pid)
			if err != nil {
				return fmt.Errorf("failed to create process handle: %v", err)
			}

			// 获取进程名称（用于日志）
			procName, err := proc.Name()
			if err != nil {
				procName = "unknown"
			}

			// 终止进程
			if err := proc.Kill(); err != nil {
				return fmt.Errorf("failed to kill process %s (PID: %d): %v", procName, conn.Pid, err)
			}

			return nil
		}
	}

	return nil // 没有进程占用该端口
}
