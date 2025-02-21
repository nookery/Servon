package routers

import (
	"servon/core/internal/managers"
	"servon/core/internal/web/controllers"

	"github.com/gin-gonic/gin"
)

func SetupDeployRouter(router *gin.RouterGroup, manager *managers.FullManager) {
	deployController := controllers.NewDeployController(manager)

	deployRouter := router.Group("/deploy")
	deployRouter.POST("/repository", deployController.DeployRepository)
}
