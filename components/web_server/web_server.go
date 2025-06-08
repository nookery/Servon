// Package web_server 提供Web服务器功能
package web_server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sevlyar/go-daemon"
)

// WebServer 封装服务器相关功能
type WebServer struct {
	*gin.Engine
	config WebServerConfig
	server *http.Server
}

// NewWebServer 创建新的Web服务器实例
func NewWebServer(config WebServerConfig) *WebServer {
	// 设置为发布模式以禁用调试日志
	gin.SetMode(gin.ReleaseMode)
	return &WebServer{
		Engine: gin.Default(),
		config: config,
	}
}

// Start 启动服务器
func (ws *WebServer) Start() error {
	addr := fmt.Sprintf("%s:%d", ws.config.Host, ws.config.Port)
	ws.server = &http.Server{
		Addr:    addr,
		Handler: ws.Engine,
	}

	return ws.server.ListenAndServe()
}

// StartWithGracefulShutdown 启动服务器并支持优雅关闭
func (ws *WebServer) StartWithGracefulShutdown() error {
	addr := fmt.Sprintf("%s:%d", ws.config.Host, ws.config.Port)
	ws.server = &http.Server{
		Addr:    addr,
		Handler: ws.Engine,
	}

	// 启动服务器
	go func() {
		if err := ws.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("服务器启动失败: %v\n", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return ws.server.Shutdown(ctx)
}

// Stop 停止服务器
func (ws *WebServer) Stop() error {
	if ws.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return ws.server.Shutdown(ctx)
	}
	return nil
}

// GetConfig 获取服务器配置
func (ws *WebServer) GetConfig() WebServerConfig {
	return ws.config
}

// SetConfig 设置服务器配置
func (ws *WebServer) SetConfig(config WebServerConfig) {
	ws.config = config
}

// IsPortAvailable 检查端口是否可用
func IsPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// FindAvailablePort 查找可用端口
func FindAvailablePort(startPort int) int {
	for port := startPort; port < startPort+100; port++ {
		if IsPortAvailable(port) {
			return port
		}
	}
	return -1
}

// StartDaemon 以守护进程方式启动
func (ws *WebServer) StartDaemon() error {
	cntxt := &daemon.Context{
		PidFileName: "webserver.pid",
		PidFilePerm: 0644,
		LogFileName: "webserver.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

	d, err := cntxt.Reborn()
	if err != nil {
		return err
	}
	if d != nil {
		return nil
	}
	defer cntxt.Release()

	return ws.StartWithGracefulShutdown()
}

// RedirectOutput 重定向输出到文件
func RedirectOutput(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// 重定向标准输出和标准错误到文件
	os.Stdout = file
	os.Stderr = file

	return nil
}

// GetServerInfo 获取服务器信息
func (ws *WebServer) GetServerInfo() map[string]interface{} {
	return map[string]interface{}{
		"host": ws.config.Host,
		"port": ws.config.Port,
		"url":  fmt.Sprintf("http://%s:%d", ws.config.Host, ws.config.Port),
	}
}

// SetupCORS 设置CORS
func (ws *WebServer) SetupCORS() {
	ws.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// SetupLogging 设置日志中间件
func (ws *WebServer) SetupLogging() {
	ws.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
}

// HealthCheck 健康检查端点
func (ws *WebServer) SetupHealthCheck() {
	ws.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})
}

// SetupStaticFiles 设置静态文件服务
func (ws *WebServer) SetupStaticFiles(relativePath, root string) {
	ws.Static(relativePath, root)
}

// SetupTemplates 设置模板
func (ws *WebServer) SetupTemplates(pattern string) {
	ws.LoadHTMLGlob(pattern)
}

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

// SetPort 设置端口号
func (ws *WebServer) SetPort(port int) {
	ws.config.Port = port
}

// SetHost 设置主机
func (ws *WebServer) SetHost(host string) {
	ws.config.Host = host
}

// RunInBackground 在后台运行服务器（作为独立进程）
func (ws *WebServer) RunInBackground() error {
	// 检查服务器是否已经在运行
	if pid, err := findProcessByPort(ws.config.Port); err == nil && pid > 0 {
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
	return ws.Stop()
}

// StopBackground 停止后台运行的服务器
func (ws *WebServer) StopBackground() error {
	// 先通过端口检查进程
	if pid, err := findProcessByPort(ws.config.Port); err == nil && pid > 0 {
		// 发送终止信号
		process, err := os.FindProcess(int(pid))
		if err == nil {
			process.Signal(syscall.SIGTERM)
		}
	}

	// 清理 PID 文件
	pidFile := "servon.pid"
	if err := os.Remove(pidFile); err == nil {
		fmt.Println("服务器已关闭")
	} else {
		fmt.Println("服务器未在运行")
	}

	return nil
}

// findProcessByPort 通过端口查找进程ID
func findProcessByPort(port int) (int32, error) {
	// 简单实现，实际项目中可能需要更复杂的逻辑
	cmd := exec.Command("lsof", "-ti", fmt.Sprintf(":%d", port))
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	if len(output) == 0 {
		return 0, fmt.Errorf("no process found")
	}
	pid, err := strconv.Atoi(string(output[:len(output)-1]))
	if err != nil {
		return 0, err
	}
	return int32(pid), nil
}
