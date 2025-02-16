package routers

import (
	"servon/core/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupPortRouter(r *gin.RouterGroup) {
	// 端口管理相关API
	api := r.Group("/ports")
	api.GET("/", handlers.HandlePortList)
	api.GET("", handlers.HandlePortList)
}
