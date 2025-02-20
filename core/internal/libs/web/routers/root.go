package routers

import (
	"servon/core/internal/integrations"
	"servon/core/internal/libs/managers"
	"servon/core/internal/libs/utils"

	"github.com/gin-gonic/gin"
)

var printer = utils.DefaultPrinter

func Setup(manager *managers.FullManager, fullIntegration *integrations.FullIntegration, r *gin.Engine, isDev bool) {
	api := r.Group("/web_api")

	SetupSoftRouter(api, manager)
	SetupProcessRouter(api, manager)
	SetupInfoRouter(api, manager)
	SetupGitHubRouter(api, fullIntegration.GitHubIntegration)
	SetupTaskRouter(api, manager)
	SetupCronRouter(api, manager)
	SetupFileRouter(api, manager)
	SetupPortRouter(api, manager)
	SetupLogsRouter(api, manager)
	SetupUserRouter(api, manager)
	SetupDeployRouter(api, manager)
	SetupIntegrationRouter(api, fullIntegration.GitHubIntegration)
	SetupUIRoutes(r)
}
