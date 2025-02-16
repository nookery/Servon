package routers

import (
	"servon/core/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupSoftRouter(r *gin.RouterGroup) {
	// 软件管理相关API
	api := r.Group("/soft")
	api.GET("", handlers.HandleGetSoftwareList)
	api.GET("/", handlers.HandleGetSoftwareList)
	api.POST("/:name/install", handlers.HandleInstallSoftware)
	api.POST("/:name/uninstall", handlers.HandleUninstallSoftware)
	api.POST("/:name/stop", handlers.HandleStopSoftware)
	api.POST("/:name/start", handlers.HandleStartSoftware)
	api.GET("/:name/status", handlers.HandleGetSoftwareStatus)
}
