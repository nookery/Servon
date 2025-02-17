package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleSystemResources 处理系统资源监控的请求
func (h *WebHandler) HandleSystemResources(c *gin.Context) {
	resources, err := h.App.SystemResourcesManager.GetSystemResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resources)
}

// HandleBasicInfo 处理基本系统信息的请求
func (h *WebHandler) HandleBasicInfo(c *gin.Context) {
	info, err := h.App.BasicInfoManager.GetBasicSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

// HandleCurrentUser 处理获取当前用户的请求
func (h *WebHandler) HandleCurrentUser(c *gin.Context) {
	user, err := h.App.SystemResourcesManager.GetCurrentUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": user})
}

// HandleNetworkResources 处理网络资源监控的请求
func (h *WebHandler) HandleNetworkResources(c *gin.Context) {
	networkStats, err := h.App.NetworkManager.GetNetworkResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, networkStats)
}

// HandleOSInfo 处理获取操作系统信息的请求
func (h *WebHandler) HandleOSInfo(c *gin.Context) {
	osInfo, err := h.App.OSInfoManager.GetOSInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"os_info": osInfo})
}

// HandleIPInfo 处理获取IP配置信息的请求
func (h *WebHandler) HandleIPInfo(c *gin.Context) {
	ipConfig, err := h.App.NetworkManager.GetIPConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ipConfig)
}
