package routers

import (
	"servon/core/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupHomeRouter(r *gin.RouterGroup) {
	r.GET("/", handlers.HandleHome)
}
