package managers

import (
	"servon/components/network_util"
)

var DefaultNetworkManager = NewNetworkManager()

type NetworkManager struct {
	networkUtil *network_util.NetworkUtil
}

func NewNetworkManager() *NetworkManager {
	return &NetworkManager{
		networkUtil: network_util.NewNetworkUtil(),
	}
}

// 为了保持向后兼容，重新导出类型
type NetworkStats = network_util.NetworkStats
type IPConfig = network_util.IPConfig
type LocalIPInfo = network_util.LocalIPInfo
type NetworkCard = network_util.NetworkCard

// GetNetworkResources 获取网络资源使用情况
func (p *NetworkManager) GetNetworkResources() (*NetworkStats, error) {
	return p.networkUtil.GetNetworkStats()
}

// GetIPConfig 获取完整的IP配置信息
func (p *NetworkManager) GetIPConfig() (*IPConfig, error) {
	return p.networkUtil.GetIPConfig()
}
