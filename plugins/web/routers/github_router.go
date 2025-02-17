package routers

import (
	"github.com/gin-gonic/gin"
)

func (w *WebRouter) SetupGitHubRouter(r *gin.RouterGroup) {
	group := r.Group("/github")
	group.POST("/setup", w.Handler.HandleGitHubSetup)
	group.GET("/callback", w.Handler.HandleGitHubCallback)
	group.POST("/webhook", w.Handler.HandleGitHubWebhook)
	group.GET("/webhooks", w.Handler.HandleGetWebhooks)
}
