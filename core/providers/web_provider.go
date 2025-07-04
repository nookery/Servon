package providers

import (
	"servon/components/web_server"
	"servon/core/managers"
	"servon/core/models"
	"servon/core/web/routers"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WebProvider struct {
	Server *web_server.WebServer
	config *models.WebConfig
}

func NewWebProvider(manager *managers.FullManager, host string, port int) *WebProvider {
	serverConfig := web_server.WebServerConfig{
		Host: host,
		Port: port,
	}
	server := web_server.NewWebServer(serverConfig)
	config := &models.WebConfig{
		Host: host,
		Port: port,
	}

	webProvider := &WebProvider{
		Server: server,
		config: config,
	}

	routers.Setup(manager, server.Engine, true)

	return webProvider
}

// GetPort 获取端口号
func (w *WebProvider) GetPort() int {
	return w.config.Port
}

// GetPortString 获取端口号字符串
func (w *WebProvider) GetPortString() string {
	return strconv.Itoa(w.config.Port)
}

// GetHost 获取主机
func (w *WebProvider) GetHost() string {
	return w.config.Host
}

// SetPort 设置端口号
func (w *WebProvider) SetPort(port int) {
	w.config.Port = port
}

// SetHost 设置主机
func (w *WebProvider) SetHost(host string) {
	w.config.Host = host
}

// GetRouter 获取路由
func (w *WebProvider) GetRouter() *gin.Engine {
	return w.Server.Engine
}
