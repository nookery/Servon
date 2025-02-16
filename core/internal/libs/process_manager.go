package libs

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/process"
)

var DefaultProcessManager = newProcessManager()

type ProcessManager struct{}

func newProcessManager() *ProcessManager {
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

// KillProcess 结束指定PID的进程
func (p *ProcessManager) KillProcess(pid int) error {
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return fmt.Errorf("获取进程失败: %v", err)
	}

	err = proc.Kill()
	if err != nil {
		return fmt.Errorf("结束进程失败: %v", err)
	}

	return nil
}
