package routers

import (
	"servon/core/internal/libs/managers"
	"servon/core/internal/web/controllers"

	"github.com/gin-gonic/gin"
)

func SetupPortRouter(r *gin.RouterGroup, manager *managers.FullManager) {
	controller := controllers.NewPortController(manager)

	// 端口管理相关API
	api := r.Group("/ports")
	api.GET("/", controller.HandlePortList)
	api.GET("", controller.HandlePortList)
}
