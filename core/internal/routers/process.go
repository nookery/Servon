package routers

import (
	"servon/core/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupProcessRouter(r *gin.RouterGroup) {
	// 进程管理相关API
	api := r.Group("/processes")
	api.GET("/", handlers.HandleProcessList)
	api.GET("", handlers.HandleProcessList)
}
