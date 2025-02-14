package libs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	maxRetries = 3 // 最大重试次数
)

type DownloadManager struct{}

func NewDownloadManager() *DownloadManager {
	return &DownloadManager{}
}

// Download 下载文件
// url: 下载地址
// destPath: 目标文件路径（包含文件名）
func (d *DownloadManager) Download(url string, destPath string) error {
	var lastErr error

	// 首先尝试不使用代理下载
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			PrintCommandOutput(fmt.Sprintf("下载失败，第 %d 次重试...", i))
			time.Sleep(time.Second * 2) // 重试前等待一段时间
		}

		err := d.downloadFile(url, destPath)
		if err == nil {
			return nil
		}
		lastErr = err
	}

	// 如果常规下载失败，尝试使用代理
	if !DefaultProxyManager.IsProxyOn() {
		PrintCommandOutput("常规下载失败，尝试开启代理重新下载...")
		software, err := DefaultProxyManager.OpenProxy()
		if err != nil {
			return fmt.Errorf("开启代理失败: %v，上一次下载错误: %v", err, lastErr)
		}

		PrintAlert("使用代理软件: " + software + " 下载...")

		// 使用代理重试下载
		for i := 0; i < maxRetries; i++ {
			if i > 0 {
				PrintCommandOutput(fmt.Sprintf("代理下载失败，第 %d 次重试...", i))
				time.Sleep(time.Second * 2)
			}

			err := d.downloadFile(url, destPath)
			if err == nil {
				return nil
			}
			lastErr = err
		}
	}

	return fmt.Errorf("下载失败（已尝试使用代理）: %v", lastErr)
}

// downloadFile 实际执行下载操作
func (d *DownloadManager) downloadFile(url string, destPath string) error {
	// 创建 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("创建下载请求失败: %s", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，服务器返回状态码: %d", resp.StatusCode)
	}

	// 创建目标文件
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %s", err)
	}
	defer out.Close()

	// 创建进度读取器
	totalSize := resp.ContentLength
	progress := &ProgressReader{
		Reader:      resp.Body,
		Total:       totalSize,
		minInterval: 800 * time.Millisecond, // 设置最小更新间隔为800ms
		lastUpdate:  time.Now(),
		OnProgress: func(current, total int64) {
			if total > 0 {
				percentage := float64(current) / float64(total) * 100
				currentSize := formatSize(current)
				totalSize := formatSize(total)
				PrintCommandOutput(fmt.Sprintf("下载进度: %.1f%% (%s/%s)", percentage, currentSize, totalSize))
			}
		},
	}

	// 下载文件
	_, err = io.Copy(out, progress)
	if err != nil {
		return fmt.Errorf("下载失败: %s", err)
	}
	PrintLn()

	return nil
}

// ProgressReader 用于跟踪读取进度
type ProgressReader struct {
	io.Reader
	Total       int64
	Current     int64
	OnProgress  func(current, total int64)
	lastUpdate  time.Time     // 上次更新时间
	minInterval time.Duration // 最小更新间隔
}

// Read 实现了 io.Reader 接口
func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.Current += int64(n)

	// 检查是否需要更新进度
	if pr.OnProgress != nil && time.Since(pr.lastUpdate) >= pr.minInterval {
		pr.OnProgress(pr.Current, pr.Total)
		pr.lastUpdate = time.Now()
	}
	return n, err
}

// formatSize 将字节大小转换为人类可读的格式
func formatSize(bytes int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
