package routers

import (
	"github.com/gin-gonic/gin"
)

func (w *WebRouter) SetupInfoRouter(r *gin.RouterGroup) {
	api := r.Group("/info")
	api.GET("/resources", w.Handler.HandleSystemResources)
	api.GET("/network", w.Handler.HandleNetworkResources)
	api.GET("/user", w.Handler.HandleCurrentUser)
	api.GET("/os", w.Handler.HandleOSInfo)
	api.GET("/basic", w.Handler.HandleBasicInfo)
	api.GET("/ip", w.Handler.HandleIPInfo)
}
