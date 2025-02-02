package handler

import (
	"fmt"
	"net/http"
	"servon/internal/system"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

// HandleSystemResources 处理系统资源监控的请求
func (h *Handler) HandleSystemResources(c *gin.Context) {
	resources, err := system.GetSystemResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resources)
}

// HandleBasicInfo 处理基本系统信息的请求
func (h *Handler) HandleBasicInfo(c *gin.Context) {
	info, err := system.GetBasicSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

// HandleSoftwareList 处理软件列表的请求
func (h *Handler) HandleSoftwareList(c *gin.Context) {
	names := system.GetSoftwareList()
	c.JSON(http.StatusOK, names)
}

// HandleSoftwareInstall 处理软件安装请求
func (h *Handler) HandleSoftwareInstall(c *gin.Context) {
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

// HandleSoftwareUninstall 处理软件卸载请求
func (h *Handler) HandleSoftwareUninstall(c *gin.Context) {
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

// HandleSoftwareStop 处理软件停止请求
func (h *Handler) HandleSoftwareStop(c *gin.Context) {
	name := c.Param("name")
	if err := system.StopSoftware(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "服务已停止"})
}

// HandleCurrentUser 处理获取当前用户的请求
func (h *Handler) HandleCurrentUser(c *gin.Context) {
	user, err := system.GetCurrentUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": user})
}

// HandleSoftwareStatus 处理获取软件状态的请求
func (h *Handler) HandleSoftwareStatus(c *gin.Context) {
	name := c.Param("name")
	status, err := system.GetSoftwareStatus(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}

// HandleProcessList 处理获取进程列表的请求
func (h *Handler) HandleProcessList(c *gin.Context) {
	processes, err := system.GetProcessList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, processes)
}

// HandleFileList 处理获取文件列表的请求
func (h *Handler) HandleFileList(c *gin.Context) {
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

// HandlePortList 处理获取端口列表的请求
func (h *Handler) HandlePortList(c *gin.Context) {
	ports, err := system.GetPortList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取端口列表失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ports)
}

// HandleOSInfo 处理获取操作系统信息的请求
func (h *Handler) HandleOSInfo(c *gin.Context) {
	osInfo, err := system.GetOSInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"os_info": osInfo})
}

// ... 其他 handler 方法 ...
