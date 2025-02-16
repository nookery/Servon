package handlers

import (
	"net/http"
	"servon/core/internal/libs"

	"github.com/gin-gonic/gin"
)

var SystemResourcesManager = libs.DefaultSystemResourcesManager
var BasicInfoManager = libs.DefaultBasicInfoManager
var NetworkManager = libs.DefaultNetworkManager
var OSInfoManager = libs.DefaultOSInfoManager

// HandleSystemResources 处理系统资源监控的请求
func HandleSystemResources(c *gin.Context) {
	resources, err := SystemResourcesManager.GetSystemResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resources)
}

// HandleBasicInfo 处理基本系统信息的请求
func HandleBasicInfo(c *gin.Context) {
	info, err := BasicInfoManager.GetBasicSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

// HandleCurrentUser 处理获取当前用户的请求
func HandleCurrentUser(c *gin.Context) {
	user, err := SystemResourcesManager.GetCurrentUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": user})
}

// HandleNetworkResources 处理网络资源监控的请求
func HandleNetworkResources(c *gin.Context) {
	networkStats, err := NetworkManager.GetNetworkResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, networkStats)
}

// HandleOSInfo 处理获取操作系统信息的请求
func HandleOSInfo(c *gin.Context) {
	osInfo, err := OSInfoManager.GetOSInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"os_info": osInfo})
}

// HandleIPInfo 处理获取IP配置信息的请求
func HandleIPInfo(c *gin.Context) {
	ipConfig, err := NetworkManager.GetIPConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ipConfig)
}
