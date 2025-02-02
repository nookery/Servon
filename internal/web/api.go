package web

import (
	"fmt"
	"net/http"
	"servon/internal/system"

	"github.com/gin-gonic/gin"
)

// setupAPIRoutes 设置所有API路由
func (s *Server) setupAPIRoutes() {
	api := s.router.Group("/web_api")
	{
		api.GET("/system/resources", s.handleSystemResources)
		api.GET("/system/user", s.handleCurrentUser)
		api.GET("/system/basic", s.handleBasicInfo)
		api.GET("/system/software", s.handleSoftwareList)
		api.GET("/system/software/:name/install", s.handleSoftwareInstall)
		api.GET("/system/software/:name/uninstall", s.handleSoftwareUninstall)
		api.POST("/system/software/:name/stop", s.handleSoftwareStop)
		api.GET("/system/software/:name/status", s.handleSoftwareStatus)
		api.GET("/system/processes", s.handleProcessList)
		api.GET("/system/files", s.handleFileList)
		api.GET("/system/ports", s.handlePortList)
	}
}

// handleSystemResources 处理系统资源监控的请求
func (s *Server) handleSystemResources(c *gin.Context) {
	resources, err := system.GetSystemResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resources)
}

// handleBasicInfo 处理基本系统信息的请求
func (s *Server) handleBasicInfo(c *gin.Context) {
	info, err := system.GetBasicSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

// handleSoftwareList 处理软件列表的请求
func (s *Server) handleSoftwareList(c *gin.Context) {
	names := system.GetSoftwareList()
	c.JSON(http.StatusOK, names)
}

// handleSoftwareInstall 处理软件安装请求
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

// handleSoftwareUninstall 处理软件卸载请求
func (s *Server) handleSoftwareUninstall(c *gin.Context) {
	name := c.Param("name")

	// 设置 SSE 头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// 获取输出通道
	outputChan, err := system.UninstallSoftware(name)
	if err != nil {
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

// handleSoftwareStop 处理软件停止请求
func (s *Server) handleSoftwareStop(c *gin.Context) {
	name := c.Param("name")
	if err := system.StopSoftware(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "服务已停止"})
}

// handleCurrentUser 处理获取当前用户的请求
func (s *Server) handleCurrentUser(c *gin.Context) {
	user, err := system.GetCurrentUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": user})
}

// handleSoftwareStatus 处理获取软件状态的请求
func (s *Server) handleSoftwareStatus(c *gin.Context) {
	name := c.Param("name")
	status, err := system.GetSoftwareStatus(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}

// handleProcessList 处理获取进程列表的请求
func (s *Server) handleProcessList(c *gin.Context) {
	processes, err := system.GetProcessList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, processes)
}

// handleFileList 处理获取文件列表的请求
func (s *Server) handleFileList(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		path = "/"
	}

	files, err := system.GetFileList(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取文件列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, files)
}

// handlePortList 处理获取端口列表的请求
func (s *Server) handlePortList(c *gin.Context) {
	ports, err := system.GetPortList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取端口列表失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ports)
}
