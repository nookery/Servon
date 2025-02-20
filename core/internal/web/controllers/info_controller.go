package controllers

import (
	"net/http"
	"servon/core/internal/managers"

	"github.com/gin-gonic/gin"
)

type InfoController struct {
	*managers.FullManager
}

func NewInfoController(manager *managers.FullManager) *InfoController {
	return &InfoController{FullManager: manager}
}

// HandleSystemResources 处理系统资源监控的请求
func (h *InfoController) HandleSystemResources(c *gin.Context) {
	resources, err := h.GetSystemResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resources)
}

// HandleBasicInfo 处理基本系统信息的请求
func (h *InfoController) HandleBasicInfo(c *gin.Context) {
	info, err := h.GetBasicSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

// HandleCurrentUser 处理获取当前用户的请求
func (h *InfoController) HandleCurrentUser(c *gin.Context) {
	user, err := h.GetCurrentUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": user})
}

// HandleNetworkResources 处理网络资源监控的请求
func (h *InfoController) HandleNetworkResources(c *gin.Context) {
	networkStats, err := h.GetNetworkResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, networkStats)
}

// HandleOSInfo 处理获取操作系统信息的请求
func (h *InfoController) HandleOSInfo(c *gin.Context) {
	osInfo, err := h.GetOSInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"os_info": osInfo})
}

// HandleIPInfo 处理获取IP配置信息的请求
func (h *InfoController) HandleIPInfo(c *gin.Context) {
	ipConfig, err := h.GetIPConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ipConfig)
}
