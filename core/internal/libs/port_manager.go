package libs

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type PortManager struct{}

func NewPortManager() *PortManager {
	return &PortManager{}
}

// PortInfo 表示端口占用信息
type PortInfo struct {
	Port      int    `json:"port"`      // 端口号
	Protocol  string `json:"protocol"`  // 协议 (TCP/UDP)
	State     string `json:"state"`     // 状态
	PID       int    `json:"pid"`       // 进程ID
	Process   string `json:"process"`   // 进程名称
	Command   string `json:"command"`   // 完整命令
	User      string `json:"user"`      // 用户
	IPAddress string `json:"ipAddress"` // 监听地址
}

// GetPortList 获取系统端口占用列表
func (p *PortManager) GetPortList() ([]PortInfo, error) {
	connections, err := net.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("获取网络连接信息失败: %v", err)
	}

	var ports []PortInfo
	for _, conn := range connections {
		// 只获取 LISTEN 状态的连接
		if conn.Status != "LISTEN" {
			continue
		}

		var processName string
		var command string
		var user string

		if conn.Pid > 0 {
			proc, err := process.NewProcess(conn.Pid)
			if err == nil {
				if name, err := proc.Name(); err == nil {
					processName = name
				}
				if cmd, err := proc.Cmdline(); err == nil {
					command = cmd
				}
				if username, err := proc.Username(); err == nil {
					user = username
				}
			}
		}

		ports = append(ports, PortInfo{
			Port:      int(conn.Laddr.Port),
			Protocol:  protocolToString(conn.Type),
			State:     conn.Status,
			PID:       int(conn.Pid),
			Process:   processName,
			Command:   command,
			User:      user,
			IPAddress: conn.Laddr.IP,
		})
	}

	return ports, nil
}

func protocolToString(proto uint32) string {
	switch proto {
	case 6:
		return "TCP"
	case 17:
		return "UDP"
	default:
		return "UNKNOWN"
	}
}
