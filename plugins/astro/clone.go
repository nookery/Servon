package astro

import (
	"fmt"
	"os"
	"strings"
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
	if err := a.RunShell("git", "config", "--global", "--get", "http.proxy"); err == nil {
		return err
	}
	if err := a.RunShell("git", "config", "--global", "--get", "https.proxy"); err == nil {
		return err
	}

	// 使用重试机制执行git clone
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := a.RunShell("git", "clone", address, savePath)
		if err == nil {
			a.PrintSuccess("Astro项目克隆成功")
			return nil
		}

		errMsg := err.Error()

		// 分析错误类型并提供建议
		var suggestion string
		switch {
		case contains(errMsg, "SSL_read") || contains(errMsg, "unexpected eof"):
			suggestion = "建议：这可能是网络不稳定或SSL证书问题，请检查网络连接或尝试设置 GIT_SSL_NO_VERIFY=true"
		case contains(errMsg, "RPC failed") || contains(errMsg, "HTTP/2 stream"):
			suggestion = "建议：这可能是由于网络不稳定或代理问题，可以尝试：\n" +
				"1. 设置较小的缓存大小: git config --global http.postBuffer 524288000\n" +
				"2. 禁用HTTP/2: git config --global http.version HTTP/1.1"
		case contains(errMsg, "Authentication failed"):
			suggestion = "建议：认证失败，请检查是否有仓库访问权限，或确认凭据是否正确"
		default:
			suggestion = "建议：请检查网络连接和仓库地址是否正确"
		}

		if attempt < maxRetries {
			a.PrintWarnf("克隆失败 (尝试 %d/%d): %v", attempt, maxRetries, errMsg)
			a.PrintWarnf("建议：%s", suggestion)
			time.Sleep(time.Second * retryInterval)

			// 如果目录已创建但克隆失败，删除它以便重试
			if _, err := os.Stat(savePath); err == nil {
				os.RemoveAll(savePath)
			}
		} else {
			a.PrintWarn(suggestion)
			return fmt.Errorf("克隆失败，已重试 %d 次: %v", maxRetries, errMsg)
		}
	}

	return nil
}

// 辅助函数：检查字符串是否包含特定子串
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
