package routers

import (
	"servon/core/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupLogsRouter(r *gin.RouterGroup) {
	group := r.Group("/logs")
	group.GET("/:channel", handlers.HandleStreamLogs)
}
