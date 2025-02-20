package utils

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// WebServerConfig 服务器配置
type WebServerConfig struct {
	Host string
	Port int
}

// WebServer 封装服务器相关功能
type WebServer struct {
	*gin.Engine
	config WebServerConfig
	server *http.Server
}

// NewWebServer 创建一个新的 Web 服务器实例
func NewWebServer() *WebServer {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	router.HandleMethodNotAllowed = false
	router.SetTrustedProxies(nil)

	// 配置日志
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/favicon.ico"},
		Output:    io.Discard,
	}))
	router.Use(gin.Recovery())

	// 配置 CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	return &WebServer{
		Engine: router,
	}
}

// Configure 配置服务器
func (ws *WebServer) Configure(host string, port int) {
	ws.config = WebServerConfig{
		Host: host,
		Port: port,
	}
}

// Start 启动服务器
func (ws *WebServer) Start() error {
	addr := fmt.Sprintf("%s:%d", ws.config.Host, ws.config.Port)

	printer.PrintInfof("启动 Web 服务器: http://%s", addr)

	// 自动停止占用端口的进程
	if err := DefaultProcessUtil.AutoStopPortProcess(ws.config.Port); err != nil {
		return err
	}

	ws.server = &http.Server{
		Addr:    addr,
		Handler: ws.Engine,
	}

	// 在新的 goroutine 中启动服务器
	go func() {
		if err := ws.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			printer.PrintErrorf("服务器运行错误: %v", err)
		}
	}()

	// 等待服务器真正启动
	return ws.waitForServerReady(addr)
}

// Stop 停止服务器
func (ws *WebServer) Stop(ctx context.Context) error {
	if ws.server != nil {
		return ws.server.Shutdown(ctx)
	}
	return nil
}

// waitForServerReady 等待服务器准备就绪
func (ws *WebServer) waitForServerReady(addr string) error {
	// 尝试连接服务器，最多等待 5 秒
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("服务器启动超时")
		case <-ticker.C:
			conn, err := net.DialTimeout("tcp", addr, time.Second)
			if err == nil {
				conn.Close()
				return nil
			}
		}
	}
}

// GetPort 获取端口号
func (ws *WebServer) GetPort() int {
	return ws.config.Port
}

// GetPortString 获取端口号字符串
func (ws *WebServer) GetPortString() string {
	return strconv.Itoa(ws.config.Port)
}

// SetPort 设置端口号
func (ws *WebServer) SetPort(port int) {
	ws.config.Port = port
}

// GetHost 获取主机
func (ws *WebServer) GetHost() string {
	return ws.config.Host
}

// SetHost 设置主机
func (ws *WebServer) SetHost(host string) {
	ws.config.Host = host
}

// GetRouter 获取路由
func (ws *WebServer) GetRouter() *gin.Engine {
	return ws.Engine
}

// StartWebServer 启动 Web 服务器
func (ws *WebServer) StartWebServer() error {
	return ws.Start()
}
