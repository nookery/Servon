package routers

import (
	"servon/core/internal/libs/managers"
	"servon/core/internal/libs/web/controllers"

	"github.com/gin-gonic/gin"
)

func SetupInfoRouter(r *gin.RouterGroup, manager *managers.FullManager) {
	controller := controllers.NewInfoController(manager)

	api := r.Group("/info")
	api.GET("/resources", controller.HandleSystemResources)
	api.GET("/network", controller.HandleNetworkResources)
	api.GET("/user", controller.HandleCurrentUser)
	api.GET("/os", controller.HandleOSInfo)
	api.GET("/basic", controller.HandleBasicInfo)
	api.GET("/ip", controller.HandleIPInfo)
}
