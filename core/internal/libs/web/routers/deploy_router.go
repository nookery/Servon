package routers

import (
	"servon/core/internal/libs/managers"
	"servon/core/internal/libs/web/controllers"

	"github.com/gin-gonic/gin"
)

func SetupDeployRouter(router *gin.RouterGroup, manager *managers.FullManager) {
	deployController := controllers.NewDeployController(manager)

	deployRouter := router.Group("/deploy")
	deployRouter.GET("/logs", deployController.HandleListDeployLogs)
	deployRouter.GET("/log", deployController.HandleGetDeployLog)
	deployRouter.DELETE("/log", deployController.HandleDeleteDeployLog)
}
