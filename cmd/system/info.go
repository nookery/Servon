package system

import (
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemInfo 包含系统信息
type SystemInfo struct {
	Hostname    string    `json:"hostname"`
	Platform    string    `json:"platform"`
	CPUUsage    float64   `json:"cpu_usage"`
	MemoryTotal uint64    `json:"memory_total"`
	MemoryUsed  uint64    `json:"memory_used"`
	DiskTotal   uint64    `json:"disk_total"`
	DiskUsed    uint64    `json:"disk_used"`
	Uptime      uint64    `json:"uptime"`
	LastUpdate  time.Time `json:"last_update"`
}

// GetSystemInfo 获取系统信息
func GetSystemInfo() (*SystemInfo, error) {
	info := &SystemInfo{
		LastUpdate: time.Now(),
	}

	// 获取主机信息
	hostInfo, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get host info: %v", err)
	}
	info.Hostname = hostInfo.Hostname
	info.Platform = fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion)
	info.Uptime = hostInfo.Uptime

	// 获取CPU使用率
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %v", err)
	}
	if len(cpuPercent) > 0 {
		info.CPUUsage = cpuPercent[0]
	}

	// 获取内存信息
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %v", err)
	}
	info.MemoryTotal = memInfo.Total
	info.MemoryUsed = memInfo.Used

	// 获取磁盘信息
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk partitions: %v", err)
	}

	var totalSize, usedSize uint64
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}
		totalSize += usage.Total
		usedSize += usage.Used
	}
	info.DiskTotal = totalSize
	info.DiskUsed = usedSize

	return info, nil
}

// GetBasicInfo 获取基本系统信息
func GetBasicInfo() map[string]string {
	return map[string]string{
		"OS":           runtime.GOOS,
		"Architecture": runtime.GOARCH,
		"CPUs":         fmt.Sprintf("%d", runtime.NumCPU()),
		"GoVersion":    runtime.Version(),
	}
}
