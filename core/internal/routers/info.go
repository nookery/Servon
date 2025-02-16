package routers

import (
	"servon/core/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupInfoRouter(r *gin.RouterGroup) {
	api := r.Group("/info")
	api.GET("/resources", handlers.HandleSystemResources)
	api.GET("/network", handlers.HandleNetworkResources)
	api.GET("/user", handlers.HandleCurrentUser)
	api.GET("/os", handlers.HandleOSInfo)
	api.GET("/basic", handlers.HandleBasicInfo)
}
