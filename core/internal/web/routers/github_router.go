package routers

import (
	"servon/core/internal/managers"
	"servon/core/internal/web/controllers"

	"github.com/gin-gonic/gin"
)

func SetupGitHubRouter(r *gin.RouterGroup, fullIntegration *managers.FullManager) {
	controller := controllers.NewGitHubController(fullIntegration)

	group := r.Group("/github")
	group.POST("/setup", controller.HandleGitHubSetup)
	group.GET("/callback", controller.HandleGitHubCallback)
	group.POST("/webhook", controller.HandleGitHubWebhook)
	group.GET("/webhooks", controller.HandleGetWebhooks)
}
