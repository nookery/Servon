package web

import (
	"servon/core/internal/models"
	"servon/core/internal/routers"
	"servon/core/internal/utils"

	"github.com/gin-gonic/gin"
)

type WebServerManager struct {
	router *gin.Engine
	config *models.WebConfig
}

func NewWebServerManager(host string, port int) *WebServerManager {
	router := utils.NewWebServer()
	config := &models.WebConfig{
		Host: host,
		Port: port,
	}

	server := &WebServerManager{
		router: router,
		config: config,
	}

	api := router.Group("/web_api")
	routers.SetupSoftRouter(api)
	routers.SetupProcessRouter(api)
	routers.SetupInfoRouter(api)
	routers.SetupGitHubRouter(api)
	routers.SetupTaskRouter(api)
	routers.SetupCronRouter(api)
	routers.SetupFileRouter(api)
	routers.SetupPortRouter(api)
	routers.SetupLogsRouter(api)
	routers.SetupUserRouter(api)
	routers.SetupHomeRouter(router.Group("/"))
	routers.SetupWebApiRouter(api)

	return server
}

func (w *WebServerManager) GetRouter() *gin.Engine {
	return w.router
}
