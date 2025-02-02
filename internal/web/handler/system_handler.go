package handler

import (
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

// HandleCurrentUser 处理获取当前用户的请求
func (h *Handler) HandleCurrentUser(c *gin.Context) {
	user, err := system.GetCurrentUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": user})
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
