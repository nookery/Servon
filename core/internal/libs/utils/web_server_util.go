package utils

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sevlyar/go-daemon"
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
				printer.PrintSuccess("服务器启动成功")
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

// RunUntilSignal 运行服务器直到收到终止信号
func (ws *WebServer) RunUntilSignal() error {
	if err := ws.Start(); err != nil {
		return err
	}

	// 创建信号通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	sig := <-quit
	printer.PrintInfof("收到信号 %v, 正在关闭服务器...", sig)

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := ws.Stop(ctx); err != nil {
		printer.PrintErrorf("服务器关闭错误: %v", err)
		return err
	}

	printer.PrintSuccess("服务器已关闭")
	return nil
}

// RunInBackground 在后台运行服务器（作为独立进程）
func (ws *WebServer) RunInBackground() error {
	printer.PrintInfof("在后台运行服务器: http://%s:%d", ws.config.Host, ws.config.Port)

	// 检查服务器是否已经在运行
	if pid, err := DefaultProcessUtil.FindProcessByPort(ws.config.Port); err == nil && pid > 0 {
		return fmt.Errorf("服务器已在运行中 (PID: %d)\n提示：如需重启，请使用 'servon serve restart' 命令", pid)
	}

	// 设置守护进程的上下文
	cntxt := &daemon.Context{
		PidFileName: "servon.pid",
		PidFilePerm: 0644,
		LogFileName: "servon.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

	// 启动守护进程
	d, err := cntxt.Reborn()
	if err != nil {
		return fmt.Errorf("创建守护进程失败: %v", err)
	}
	if d != nil {
		return nil // 父进程退出
	}

	// 子进程继续执行
	defer cntxt.Release()

	// 启动服务器
	if err := ws.Start(); err != nil {
		return err
	}

	// 等待信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ws.Stop(ctx)
}

// StopBackground 停止后台运行的服务器
func (ws *WebServer) StopBackground() error {
	// 先通过端口检查进程
	if pid, err := DefaultProcessUtil.FindProcessByPort(ws.config.Port); err == nil && pid > 0 {
		// 发送终止信号
		process, err := os.FindProcess(int(pid))
		if err == nil {
			process.Signal(syscall.SIGTERM)
		}
	}

	// 清理 PID 文件
	pidFile := "servon.pid"
	if err := os.Remove(pidFile); err == nil {
		printer.PrintSuccess("服务器已关闭")
	} else {
		printer.PrintSuccess("服务器未在运行")
	}

	return nil
}
