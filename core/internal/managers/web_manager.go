package managers

import (
	"fmt"
	"os"
	"servon/core/internal/models"
	"servon/core/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WebServerManager struct {
	router *gin.Engine
	config *models.WebConfig
}

func NewWebServerManager(host string, port int) *WebServerManager {
	router := utils.NewWebServer()
	config := &models.WebConfig{
		Host: host,
		Port: port,
	}

	server := &WebServerManager{
		router: router,
		config: config,
	}

	return server
}

// GetPort
func (w *WebServerManager) GetPort() int {
	return w.config.Port
}

// GetPortString
func (w *WebServerManager) GetPortString() string {
	return strconv.Itoa(w.config.Port)
}

// GetHost
func (w *WebServerManager) GetHost() string {
	return w.config.Host
}

// SetPort
func (w *WebServerManager) SetPort(port int) {
	w.config.Port = port
}

// SetHost
func (w *WebServerManager) SetHost(host string) {
	w.config.Host = host
}

func (w *WebServerManager) GetRouter() *gin.Engine {
	return w.router
}

func (p *WebServerManager) StartWebServer() {
	router := p.router

	// 启动 Web 服务器
	// p.PrintInfof("启动 Web 服务器: http://%s:%d", host, port)
	err := router.Run(fmt.Sprintf("%s:%d", p.config.Host, p.config.Port))
	if err != nil {
		printer.PrintErrorf("启动 Web 服务器失败: %v", err)
		os.Exit(1)
	}
}
