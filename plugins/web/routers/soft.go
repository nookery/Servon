package routers

import (
	"github.com/gin-gonic/gin"
)

func (w *WebRouter) SetupSoftRouter(r *gin.RouterGroup) {
	w.Handler.App.PrintInfof("Setup Soft Router")

	api := r.Group("/soft")
	api.GET("", w.Handler.HandleGetSoftwareList)
	api.GET("/", w.Handler.HandleGetSoftwareList)
	api.POST("/:name/install", w.Handler.HandleInstallSoftware)
	api.POST("/:name/uninstall", w.Handler.HandleUninstallSoftware)
	api.POST("/:name/stop", w.Handler.HandleStopSoftware)
	api.POST("/:name/start", w.Handler.HandleStartSoftware)
	api.GET("/:name/status", w.Handler.HandleGetSoftwareStatus)
}
