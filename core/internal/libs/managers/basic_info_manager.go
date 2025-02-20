package managers

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

type BasicInfoManager struct {
	BasicInfo *BasicInfo
}

func NewBasicInfoManager() *BasicInfoManager {
	return &BasicInfoManager{}
}

type BasicInfo struct {
	Hostname     string    `json:"hostname"`
	Platform     string    `json:"platform"`
	Architecture string    `json:"architecture"`
	GoVersion    string    `json:"goVersion"`
	NumCPU       int       `json:"numCPU"`
	Uptime       string    `json:"uptime"`
	StartTime    time.Time `json:"startTime"`
}

var startTime = time.Now()

func (p *BasicInfoManager) GetBasicSystemInfo() (*BasicInfo, error) {
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
		Uptime:       time.Since(startTime).Round(time.Second).String(),
		StartTime:    startTime,
	}

	return info, nil
}
