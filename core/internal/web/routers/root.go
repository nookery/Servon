package routers

import (
	"servon/core/internal/managers"

	"github.com/gin-gonic/gin"
)

func Setup(manager *managers.FullManager, r *gin.Engine, isDev bool) {
	api := r.Group("/web_api")

	SetupSoftRouter(api, manager)
	SetupProcessRouter(api, manager)
	SetupInfoRouter(api, manager)
	SetupGitHubRouter(api, manager)
	SetupTaskRouter(api, manager)
	SetupCronRouter(api, manager)
	SetupFileRouter(api, manager)
	SetupPortRouter(api, manager)
	SetupUserRouter(api, manager)
	SetupDeployRouter(api, manager)
	SetupIntegrationRouter(api, manager)
	SetupLogRouter(api, manager.LogManager)
	SetupUIRoutes(r)
}
