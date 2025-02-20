package routers

import (
	"servon/core/internal/managers"
	"servon/core/internal/web/controllers"

	"github.com/gin-gonic/gin"
)

func SetupProcessRouter(r *gin.RouterGroup, manager *managers.FullManager) {
	controller := controllers.NewProcessController(manager)

	// 进程管理相关API
	api := r.Group("/processes")
	api.GET("/", controller.HandleProcessList)
	api.GET("", controller.HandleProcessList)
}
