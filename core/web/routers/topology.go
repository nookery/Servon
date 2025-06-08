package routers

import (
	"servon/core/managers"
	"servon/core/web/controllers"

	"github.com/gin-gonic/gin"
)

func SetupTopologyRoutes(r *gin.RouterGroup, manager *managers.ProjectManager) {
	controller := controllers.NewTopologyController(manager)

	topology := r.Group("/topology")
	{
		topology.GET("/gateways", controller.HandleGetGateways)
		topology.GET("/gateways/:gateway/projects", controller.HandleGetProjects)
		topology.POST("/gateways/:gateway/projects", controller.HandleAddProject)
		topology.DELETE("/gateways/:gateway/projects/:name", controller.HandleRemoveProject)
		topology.GET("/gateways/:gateway/config", controller.HandleGetGatewayConfig)
		topology.PUT("/gateways/:gateway/config", controller.HandleSetGatewayConfig)
	}
}
