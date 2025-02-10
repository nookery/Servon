package libs

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Process struct {
	PID     int     `json:"pid"`
	User    string  `json:"user"`
	CPU     float64 `json:"cpu"`
	Memory  float64 `json:"memory"`
	Command string  `json:"command"`
}

// GetProcessList 返回系统当前运行的进程列表
func GetProcessList() ([]Process, error) {
	// 使用 ps 命令获取进程信息
	cmd := exec.Command("ps", "aux", "--no-headers")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("获取进程列表失败: %v", err)
	}

	processes := []Process{}
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 11 {
			continue
		}

		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}

		cpu, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			cpu = 0
		}

		mem, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			mem = 0
		}

		process := Process{
			PID:     pid,
			User:    fields[0],
			CPU:     cpu,
			Memory:  mem,
			Command: strings.Join(fields[10:], " "),
		}

		processes = append(processes, process)
	}

	return processes, nil
}
