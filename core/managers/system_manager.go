package managers

import (
	"fmt"
	"os/user"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemResourcesManager struct {
	SystemResources *SystemResources
}

func NewSystemResourcesManager() *SystemResourcesManager {
	return &SystemResourcesManager{}
}

type SystemResources struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
}

// GetCurrentUser 获取当前系统用户名
func (p *SystemResourcesManager) GetCurrentUser() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("获取当前用户失败: %v", err)
	}
	return currentUser.Username, nil
}

// GetSystemResources 获取系统CPU和内存使用率
func (p *SystemResourcesManager) GetSystemResources() (*SystemResources, error) {
	// 获取CPU使用率
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	// 获取内存使用情况
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// 获取根目录磁盘使用情况
	diskInfo, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	return &SystemResources{
		CPUUsage:    cpuPercent[0],        // CPU使用率
		MemoryUsage: memInfo.UsedPercent,  // 内存使用率
		DiskUsage:   diskInfo.UsedPercent, // 磁盘使用率
	}, nil
}
