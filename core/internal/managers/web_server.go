package managers

import (
	"servon/core/internal/utils"

	"github.com/gin-gonic/gin"
)

type WebServerManager struct {
	router *gin.Engine
}

func NewWebServerManager(host string, port int, withUI bool) *WebServerManager {
	return &WebServerManager{
		router: utils.NewWebServer(host, port, withUI),
	}
}

func (w *WebServerManager) GetRouter() *gin.Engine {
	return w.router
}
