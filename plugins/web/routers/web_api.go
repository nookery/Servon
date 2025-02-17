package routers

import (
	"github.com/gin-gonic/gin"
)

func (w *WebRouter) SetupWebApiRouter(r *gin.RouterGroup) {
	r.GET("/", w.Handler.HandleWebApi)
}
