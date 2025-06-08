package managers

import (
	"servon/components/system_info"
	"time"
)

type BasicInfoManager struct {
	provider *system_info.BasicInfoProvider
}

func NewBasicInfoManager() *BasicInfoManager {
	return &BasicInfoManager{
		provider: system_info.NewBasicInfoProvider(),
	}
}

// BasicInfo 为了保持向后兼容，重新导出类型
type BasicInfo = system_info.BasicInfo

// GetBasicSystemInfo 获取基本系统信息
func (p *BasicInfoManager) GetBasicSystemInfo() (*BasicInfo, error) {
	return p.provider.GetBasicSystemInfo()
}

// GetUptime 获取运行时间
func (p *BasicInfoManager) GetUptime() time.Duration {
	return p.provider.GetUptime()
}

// GetStartTime 获取启动时间
func (p *BasicInfoManager) GetStartTime() time.Time {
	return p.provider.GetStartTime()
}
