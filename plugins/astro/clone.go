package astro

import (
	"fmt"
	"os"
	"time"
)

const (
	maxRetries    = 3 // 最大重试次数
	retryInterval = 2 // 重试间隔（秒）
)

func (a *AstroPlugin) clone(address string, savePath string) error {
	// 如果路径存在，且不为空，则不能克隆
	if _, err := os.Stat(savePath); err == nil {
		return fmt.Errorf("路径 %s 已存在，不能克隆", savePath)
	}

	// 检查并输出代理配置
	if out, err := a.RunShellWithOutput("git", "config", "--global", "--get", "http.proxy"); err == nil && out != "" {
		fmt.Printf("Git HTTP 代理: %s", out)
	}
	if out, err := a.RunShellWithOutput("git", "config", "--global", "--get", "https.proxy"); err == nil && out != "" {
		fmt.Printf("Git HTTPS 代理: %s", out)
	}

	// 使用重试机制执行git clone
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := a.RunShell("git", "clone", address, savePath); err == nil {
			a.PrintSuccess("Astro项目克隆成功")
			return nil
		} else {
			lastErr = err
			if attempt < maxRetries {
				a.PrintWarnf("克隆失败 (尝试 %d/%d)，%v 秒后重试...", attempt, maxRetries, retryInterval)
				time.Sleep(time.Second * retryInterval)

				// 如果目录已创建但克隆失败，删除它以便重试
				if _, err := os.Stat(savePath); err == nil {
					os.RemoveAll(savePath)
				}
			}
		}
	}

	return fmt.Errorf("克隆失败，已重试 %d 次: %v", maxRetries, lastErr)
}
