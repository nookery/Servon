package routers

import (
	"github.com/gin-gonic/gin"
)

func (w *WebRouter) SetupLogsRouter(r *gin.RouterGroup) {
	group := r.Group("/logs")
	group.GET("/:channel", w.Handler.HandleStreamLogs)
}
