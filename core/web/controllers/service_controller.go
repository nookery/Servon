package controllers

import (
	"fmt"
	"net/http"
	"servon/core/managers"

	"github.com/gin-gonic/gin"
)

// ServiceController 处理服务管理相关请求
type ServiceController struct {
	manager *managers.ServiceManager
}

// NewServiceController 创建服务控制器实例
func NewServiceController(manager *managers.ServiceManager) *ServiceController {
	return &ServiceController{
		manager: manager,
	}
}

// GetServiceList 获取所有服务列表
func (c *ServiceController) GetServiceList(ctx *gin.Context) {

	ctx.String(http.StatusOK, "")
}

// StartService 启动指定服务
func (c *ServiceController) StartService(ctx *gin.Context) {
	name := ctx.Param("name")

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("服务 %s 已启动", name),
	})
}

// StopService 停止指定服务
func (c *ServiceController) StopService(ctx *gin.Context) {
	name := ctx.Param("name")

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("服务 %s 已停止", name),
	})
}

// RestartService 重启指定服务
func (c *ServiceController) RestartService(ctx *gin.Context) {
	name := ctx.Param("name")

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("服务 %s 已重启", name),
	})
}

// GetServiceConfig 获取服务配置
func (c *ServiceController) GetServiceConfig(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}

// UpdateServiceConfig 更新服务配置
func (c *ServiceController) UpdateServiceConfig(ctx *gin.Context) {
	name := ctx.Param("name")

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("服务 %s 配置已更新", name),
	})
}

// GetServiceLogs 获取服务日志
func (c *ServiceController) GetServiceLogs(ctx *gin.Context) {

	ctx.String(http.StatusOK, "")
}

// GetServiceDetails 获取服务详情
func (c *ServiceController) GetServiceDetails(ctx *gin.Context) {
	name := ctx.Param("name")

	ctx.JSON(http.StatusOK, name)
}

// AddBackgroundService 添加后台服务
func (c *ServiceController) AddBackgroundService(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}

// RemoveBackgroundService 删除后台服务
func (c *ServiceController) RemoveBackgroundService(ctx *gin.Context) {
	name := ctx.Param("name")

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("后台服务 %s 已删除", name),
	})
}
