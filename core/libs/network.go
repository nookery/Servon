package libs

import (
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

type NetworkStats struct {
	DownloadSpeed int64 `json:"download_speed"` // 下载速度（字节/秒）
	UploadSpeed   int64 `json:"upload_speed"`   // 上传速度（字节/秒）
}

var (
	lastStats     []net.IOCountersStat
	lastStatsTime time.Time
)

// GetNetworkResources 获取网络资源使用情况
func GetNetworkResources() (*NetworkStats, error) {
	// 获取当前网络统计信息
	currentStats, err := net.IOCounters(false) // false表示获取所有网卡的总和
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()

	// 如果是第一次获取数据
	if lastStats == nil {
		lastStats = currentStats
		lastStatsTime = currentTime
		return &NetworkStats{
			DownloadSpeed: 0,
			UploadSpeed:   0,
		}, nil
	}

	// 计算时间差（秒）
	duration := currentTime.Sub(lastStatsTime).Seconds()

	// 计算速度
	var totalBytesRecv int64
	var totalBytesSent int64
	var lastBytesRecv int64
	var lastBytesSent int64

	for _, stat := range currentStats {
		totalBytesRecv += int64(stat.BytesRecv)
		totalBytesSent += int64(stat.BytesSent)
	}

	for _, stat := range lastStats {
		lastBytesRecv += int64(stat.BytesRecv)
		lastBytesSent += int64(stat.BytesSent)
	}

	// 计算每秒的速度
	downloadSpeed := int64(float64(totalBytesRecv-lastBytesRecv) / duration)
	uploadSpeed := int64(float64(totalBytesSent-lastBytesSent) / duration)

	// 更新上次的统计信息
	lastStats = currentStats
	lastStatsTime = currentTime

	return &NetworkStats{
		DownloadSpeed: downloadSpeed,
		UploadSpeed:   uploadSpeed,
	}, nil
}
