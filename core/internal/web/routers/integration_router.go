package routers

import (
	"servon/core/internal/managers"
	"servon/core/internal/web/controllers"

	"github.com/gin-gonic/gin"
)

func SetupIntegrationRouter(r *gin.RouterGroup, fullIntegration *managers.FullManager) {
	controller := controllers.NewIntegrationController(fullIntegration)

	api := r.Group("/integrations")
	api.GET("/github/repos", controller.HandleListGitHubRepos)
	api.GET("/github/logs", controller.HandleGetGitHubLogs)
	// api.GET("/", w.Handler.HandleIntegrationList)
	// api.GET("", w.Handler.HandleIntegrationList)
}
