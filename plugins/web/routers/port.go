package routers

import (
	"github.com/gin-gonic/gin"
)

func (w *WebRouter) SetupPortRouter(r *gin.RouterGroup) {
	// 端口管理相关API
	api := r.Group("/ports")
	api.GET("/", w.Handler.HandlePortList)
	api.GET("", w.Handler.HandlePortList)
}
