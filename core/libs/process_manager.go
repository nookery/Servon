package libs

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcessManager struct{}

func NewProcessManager() *ProcessManager {
	return &ProcessManager{}
}

type Process struct {
	PID     int     `json:"pid"`
	User    string  `json:"user"`
	CPU     float64 `json:"cpu"`
	Memory  float64 `json:"memory"`
	Command string  `json:"command"`
}

// GetProcessList 返回系统当前运行的进程列表
func (p *ProcessManager) GetProcessList() ([]Process, error) {
	processes := []Process{}

	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("获取进程列表失败: %v", err)
	}

	for _, proc := range procs {
		cpu, _ := proc.CPUPercent()
		mem, _ := proc.MemoryPercent()
		username, _ := proc.Username()
		cmdline, _ := proc.Cmdline()

		process := Process{
			PID:     int(proc.Pid),
			User:    username,
			CPU:     cpu,
			Memory:  float64(mem),
			Command: cmdline,
		}

		processes = append(processes, process)
	}

	return processes, nil
}
