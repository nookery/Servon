package libs

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

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
func GetPortList() ([]PortInfo, error) {
	// 使用 netstat 命令获取端口信息
	cmd := exec.Command("netstat", "-tulpn")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("执行 netstat 命令失败: %v", err)
	}

	// 解析输出
	lines := strings.Split(string(output), "\n")
	var ports []PortInfo

	// 跳过前两行标题
	for _, line := range lines[2:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}

		// 解析地址和端口
		localAddr := fields[3]
		addrParts := strings.Split(localAddr, ":")
		if len(addrParts) != 2 {
			continue
		}

		port, err := strconv.Atoi(addrParts[1])
		if err != nil {
			continue
		}

		// 解析 PID/进程名
		pidProgram := fields[6]
		var pid int
		var process string

		if pidProgram != "-" {
			parts := strings.Split(pidProgram, "/")
			if len(parts) == 2 {
				pid, _ = strconv.Atoi(parts[0])
				process = parts[1]
			}
		}

		// 获取进程的完整命令
		var command string
		if pid > 0 {
			if cmd, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "command=").Output(); err == nil {
				command = strings.TrimSpace(string(cmd))
			}
		}

		// 获取进程的用户
		var user string
		if pid > 0 {
			if userData, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "user=").Output(); err == nil {
				user = strings.TrimSpace(string(userData))
			}
		}

		ports = append(ports, PortInfo{
			Port:      port,
			Protocol:  strings.ToUpper(fields[0]),
			State:     fields[5],
			PID:       pid,
			Process:   process,
			Command:   command,
			User:      user,
			IPAddress: addrParts[0],
		})
	}

	return ports, nil
}
