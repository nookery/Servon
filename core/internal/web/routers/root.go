package routers

import (
	"servon/core/internal/managers"
	"servon/core/internal/utils"

	"github.com/gin-gonic/gin"
)

var printer = utils.DefaultPrinter

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
	SetupLogsRouter(api, manager)
	SetupUserRouter(api, manager)
	SetupDeployRouter(api, manager)
	SetupIntegrationRouter(api, manager)
	SetupUIRoutes(r)
}
