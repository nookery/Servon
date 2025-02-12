package serve

import (
	"fmt"
	"servon/core/libs"

	"github.com/gin-gonic/gin"
)

func StartWebServer(host string, port int, withUI bool) {
	router := libs.NewWebServer(host, port, withUI)

	// 设置API路由
	setupAPIRoutes(router)

	// 如果启用了UI，设置UI路由
	if withUI {
		setupUIRoutes(router)
	}

	router.Run(fmt.Sprintf("%s:%d", host, port))
}

// setupAPIRoutes 设置所有API路由
func setupAPIRoutes(router *gin.Engine) {
	h := NewHanlder()
	api := router.Group("/web_api")
	{
		api.GET("/system/resources", h.HandleSystemResources)
		api.GET("/system/network", h.HandleNetworkResources)
		api.GET("/system/user", h.HandleCurrentUser)
		api.GET("/system/os", h.HandleOSInfo)
		api.GET("/system/basic", h.HandleBasicInfo)
		api.GET("/system/software", h.HandleGetSoftwareList)
		api.GET("/system/software/:name/install", h.HandleInstallSoftware)
		api.GET("/system/software/:name/uninstall", h.HandleUninstallSoftware)
		api.POST("/system/software/:name/stop", h.HandleStopSoftware)
		api.GET("/system/software/:name/status", h.HandleGetSoftwareStatus)
		api.GET("/system/processes", h.HandleProcessList)
		api.GET("/system/files", h.HandleFileList)
		api.GET("/system/ports", h.HandlePortList)

		// // 定时任务相关API
		// api.GET("/cron/tasks", h.HandleListCronTasks)              // 获取所有定时任务
		// api.POST("/cron/tasks", h.HandleCreateCronTask)            // 创建定时任务
		// api.PUT("/cron/tasks/:id", h.HandleUpdateCronTask)         // 更新定时任务
		// api.DELETE("/cron/tasks/:id", h.HandleDeleteCronTask)      // 删除定时任务
		// api.POST("/cron/tasks/:id/toggle", h.HandleToggleCronTask) // 启用/禁用定时任务
	}
}
