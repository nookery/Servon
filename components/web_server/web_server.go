// Package web_server 提供Web服务器功能
package web_server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WebServer 封装服务器相关功能
type WebServer struct {
	*gin.Engine
	config WebServerConfig
	server *http.Server
	logger Logger
}

// NewWebServer 创建新的Web服务器实例
func NewWebServer(config WebServerConfig) *WebServer {
	// 设置为发布模式以禁用调试日志
	gin.SetMode(gin.ReleaseMode)

	// 如果没有提供自定义 logger，使用默认 logger
	var logger Logger
	if config.Logger != nil {
		logger = config.Logger
	} else {
		logger = NewDefaultLogger()
	}

	return &WebServer{
		Engine: gin.Default(),
		config: config,
		logger: logger,
	}
}
