package routers

import (
	"servon/core/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupGitHubRouter(r *gin.RouterGroup) {
	group := r.Group("/github")
	group.POST("/setup", handlers.HandleGitHubSetup)
	group.GET("/callback", handlers.HandleGitHubCallback)
	group.POST("/webhook", handlers.HandleGitHubWebhook)
}
