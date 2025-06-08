// Package string_util 提供字符串处理功能
package string_util

import (
	"strings"
	"time"
)

var DefaultStringUtil = &StringUtil{}

type StringUtil struct {
}

// GetEmojiForBool 获取布尔值的emoji
func (s *StringUtil) GetEmojiForBool(value bool) string {
	if value {
		return "✅"
	}
	return "❌"
}

// GetProjectNameFromString 从字符串中获取项目名称
func (s *StringUtil) GetProjectNameFromString(str string) string {
	parts := strings.Split(str, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return time.Now().Format("20060102150405")
}

// GetProjectNameFromRepoURL 从仓库URL中获取项目名称
func (s *StringUtil) GetProjectNameFromRepoURL(repoURL string) string {
	repoURL = strings.TrimSuffix(repoURL, ".git")
	repoURL = strings.TrimPrefix(repoURL, "https://")
	repoURL = strings.TrimPrefix(repoURL, "http://")
	repoURL = strings.TrimPrefix(repoURL, "git@")
	repoURL = strings.TrimPrefix(repoURL, "ssh://")
	repoURL = strings.TrimPrefix(repoURL, "git+")
	repoURL = strings.TrimPrefix(repoURL, "git+ssh://")
	repoURL = strings.TrimPrefix(repoURL, "git+http://")
	repoURL = strings.TrimPrefix(repoURL, "git+https://")

	parts := strings.Split(repoURL, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return time.Now().Format("20060102150405")
}

// GetProjectNameFromFileURL 从工作目录中获取项目名称
func (s *StringUtil) GetProjectNameFromFileURL(fileURL string) string {
	parts := strings.Split(fileURL, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return time.Now().Format("20060102150405")
}