package routers

import (
	"servon/core/internal/libs/managers"
	"servon/core/internal/web/controllers"

	"github.com/gin-gonic/gin"
)

func SetupSoftRouter(r *gin.RouterGroup, manager *managers.FullManager) {
	controller := controllers.NewSoftController(manager)

	api := r.Group("/soft")
	api.GET("", controller.HandleGetSoftwareList)
	api.GET("/", controller.HandleGetSoftwareList)
	api.POST("/:name/install", controller.HandleInstallSoftware)
	api.POST("/:name/uninstall", controller.HandleUninstallSoftware)
	api.POST("/:name/stop", controller.HandleStopSoftware)
	api.POST("/:name/start", controller.HandleStartSoftware)
	api.GET("/:name/status", controller.HandleGetSoftwareStatus)
}
