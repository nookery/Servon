package system_info

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

// BasicInfo 表示基本系统信息
type BasicInfo struct {
	Hostname     string    `json:"hostname"`
	Platform     string    `json:"platform"`
	Architecture string    `json:"architecture"`
	GoVersion    string    `json:"goVersion"`
	NumCPU       int       `json:"numCPU"`
	Uptime       string    `json:"uptime"`
	StartTime    time.Time `json:"startTime"`
}

// BasicInfoProvider 提供基本系统信息的核心功能
type BasicInfoProvider struct {
	startTime time.Time
}

// NewBasicInfoProvider 创建新的基本信息提供者
func NewBasicInfoProvider() *BasicInfoProvider {
	return &BasicInfoProvider{
		startTime: time.Now(),
	}
}

// GetBasicSystemInfo 获取基本系统信息
func (p *BasicInfoProvider) GetBasicSystemInfo() (*BasicInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("获取主机名失败: %v", err)
	}

	info := &BasicInfo{
		Hostname:     hostname,
		Platform:     runtime.GOOS,
		Architecture: runtime.GOARCH,
		GoVersion:    runtime.Version(),
		NumCPU:       runtime.NumCPU(),
		Uptime:       time.Since(p.startTime).Round(time.Second).String(),
		StartTime:    p.startTime,
	}

	return info, nil
}

// GetUptime 获取运行时间
func (p *BasicInfoProvider) GetUptime() time.Duration {
	return time.Since(p.startTime)
}

// GetStartTime 获取启动时间
func (p *BasicInfoProvider) GetStartTime() time.Time {
	return p.startTime
}
