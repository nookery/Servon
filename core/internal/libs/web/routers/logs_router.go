package routers

import (
	"servon/core/internal/libs/managers"
	"servon/core/internal/libs/web/controllers"

	"github.com/gin-gonic/gin"
)

func SetupLogsRouter(r *gin.RouterGroup, manager *managers.FullManager) {
	controller := controllers.NewLogController(manager)

	group := r.Group("/logs")
	group.GET("/:channel", controller.HandleStreamLogs)
}
