package web

import (
	"fmt"
	"net/http"

	"servon/internal/system"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	port   int
}

func NewServer(port int) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

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

	return &Server{
		router: router,
		port:   port,
	}
}

func (s *Server) setupRoutes() {
	// API 路由
	api := s.router.Group("/api")
	{
		api.GET("/system/basic", s.handleBasicInfo)                        // 新增基本信息接口
		api.GET("/system/software", s.handleSoftwareList)                  // 新增软件列表接口
		api.GET("/system/software/:name/install", s.handleSoftwareInstall) // 改为 GET 方法以支持 SSE
		api.POST("/system/software/:name/uninstall", s.handleSoftwareUninstall)
		api.POST("/system/software/:name/stop", s.handleSoftwareStop) // 新增停止服务接口
		api.GET("/system/user", s.handleCurrentUser)                  // 新增获取当前用户接口
		// TODO: 添加更多API路由
	}
}

func (s *Server) Start() error {
	s.setupRoutes()
	return s.router.Run(fmt.Sprintf(":%d", s.port))
}

// API处理函数
func (s *Server) handleSystemInfo(c *gin.Context) {
	info, err := system.GetSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

// 新增：处理基本系统信息的接口
func (s *Server) handleBasicInfo(c *gin.Context) {
	info, err := system.GetBasicSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

// 新增：处理软件列表的接口
func (s *Server) handleSoftwareList(c *gin.Context) {
	info, err := system.GetSoftwareList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

// 处理软件安装请求
func (s *Server) handleSoftwareInstall(c *gin.Context) {
	name := c.Param("name")

	// 设置 SSE 头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// 获取输出通道
	outputChan, err := system.InstallSoftware(name)
	if err != nil {
		// 发送错误消息
		c.SSEvent("message", fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	// 清空缓冲区
	if f, ok := c.Writer.(http.Flusher); ok {
		f.Flush()
	}

	// 发送输出
	for msg := range outputChan {
		c.SSEvent("message", msg)
		if f, ok := c.Writer.(http.Flusher); ok {
			f.Flush()
		}
	}
}

// 处理软件卸载请求
func (s *Server) handleSoftwareUninstall(c *gin.Context) {
	name := c.Param("name")
	if err := system.UninstallSoftware(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "卸载成功"})
}

// 处理软件服务停止请求
func (s *Server) handleSoftwareStop(c *gin.Context) {
	name := c.Param("name")
	if err := system.StopSoftware(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "服务已停止"})
}

// 处理获取当前用户的请求
func (s *Server) handleCurrentUser(c *gin.Context) {
	user, err := system.GetCurrentUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": user})
}
