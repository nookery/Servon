package routers

import (
	"servon/core/internal/libs/integrations"
	"servon/core/internal/libs/web/controllers"

	"github.com/gin-gonic/gin"
)

func SetupGitHubRouter(r *gin.RouterGroup, fullIntegration *integrations.FullIntegration) {
	controller := controllers.NewGitHubController(fullIntegration)

	group := r.Group("/github")
	group.POST("/setup", controller.HandleGitHubSetup)
	group.GET("/callback", controller.HandleGitHubCallback)
	group.POST("/webhook", controller.HandleGitHubWebhook)
	group.GET("/webhooks", controller.HandleGetWebhooks)
}
