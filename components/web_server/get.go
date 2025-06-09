package web_server

import (
	"fmt"
	"strconv"
)

// GetPort 获取端口
func (ws *WebServer) GetPort() int {
	return ws.config.Port
}

// GetHost 获取主机
func (ws *WebServer) GetHost() string {
	return ws.config.Host
}

// GetURL 获取完整URL
func (ws *WebServer) GetURL() string {
	return fmt.Sprintf("http://%s:%d", ws.config.Host, ws.config.Port)
}

// GetPortString 获取端口字符串
func (ws *WebServer) GetPortString() string {
	return strconv.Itoa(ws.config.Port)
}

// GetServerInfo 获取服务器信息
func (ws *WebServer) GetServerInfo() map[string]interface{} {
	return map[string]any{
		"host": ws.config.Host,
		"port": ws.config.Port,
		"url":  fmt.Sprintf("http://%s:%d", ws.config.Host, ws.config.Port),
	}
}

// GetConfig 获取服务器配置
func (ws *WebServer) GetConfig() WebServerConfig {
	return ws.config
}

// GetLogger 获取当前使用的日志器
func (ws *WebServer) GetLogger() Logger {
	return ws.logger
}
