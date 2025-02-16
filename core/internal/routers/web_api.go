package routers

import (
	"servon/core/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupWebApiRouter(r *gin.RouterGroup) {
	r.GET("/", handlers.HandleWebApi)
}
