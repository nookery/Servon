package routers

import (
	"github.com/gin-gonic/gin"
)

func (w *WebRouter) SetupProcessRouter(r *gin.RouterGroup) {
	// 进程管理相关API
	api := r.Group("/processes")
	api.GET("/", w.Handler.HandleProcessList)
	api.GET("", w.Handler.HandleProcessList)
}
