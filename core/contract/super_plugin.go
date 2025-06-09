package contract

import (
	"github.com/gin-gonic/gin"
)

// SuperPlugin 定义完整的插件接口，组合Web和CLI功能
type SuperPlugin interface {
	SuperWebPlugin
	SuperCLIPlugin
}

// SuperWebPlugin 定义插件Web功能接口
type SuperWebPlugin interface {
	// GetName 返回插件名称
	GetName() string
	// SetupWebRoutes 设置Web路由
	SetupWebRoutes(router *gin.RouterGroup, manager interface{})
	// GetAPIPrefix 返回API前缀路径
	GetAPIPrefix() string
	// GetDescription 返回插件描述
	GetDescription() string
	// IsEnabled 检查插件是否启用
	IsEnabled() bool
}

// SuperCLIPlugin 定义插件CLI功能接口
type SuperCLIPlugin interface {
	// GetName 返回插件名称
	GetName() string
	// GetDescription 返回插件描述
	GetDescription() string
	// IsEnabled 检查插件是否启用
	IsEnabled() bool
}

// PluginInfo 插件基本信息
type PluginInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
	Enabled     bool   `json:"enabled"`
	HasWebAPI   bool   `json:"has_web_api"`
	HasCLI      bool   `json:"has_cli"`
	APIPrefix   string `json:"api_prefix,omitempty"`
}
