package routers

import (
	"github.com/gin-gonic/gin"
)

func (w *WebRouter) SetupIntegrationRouter(r *gin.RouterGroup) {
	api := r.Group("/integrations")
	api.GET("/github/repos", w.Handler.HandleListGitHubRepos)
	api.GET("/github/logs", w.Handler.HandleGetGitHubLogs)
	// api.GET("/", w.Handler.HandleIntegrationList)
	// api.GET("", w.Handler.HandleIntegrationList)
}
