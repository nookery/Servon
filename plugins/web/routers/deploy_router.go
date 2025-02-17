package routers

import (
	"github.com/gin-gonic/gin"
)

func (w *WebRouter) SetupDeployRouter(router *gin.RouterGroup) {
	deployRouter := router.Group("/deploy")
	deployRouter.GET("/logs", w.Handler.HandleListDeployLogs)
	deployRouter.GET("/log", w.Handler.HandleGetDeployLog)
	deployRouter.DELETE("/log", w.Handler.HandleDeleteDeployLog)
}
