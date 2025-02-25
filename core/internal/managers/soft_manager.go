package managers

import (
	"servon/core/internal/managers/soft"
)

// SoftManager 是软件管理器的公共接口
// 这个文件作为适配层，保持向后兼容性
type SoftManager = soft.Manager

// NewSoftManager 创建新的软件管理器
func NewSoftManager(logDir string) *SoftManager {
	return soft.NewManager(logDir)
}
