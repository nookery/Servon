package routers

import (
	"github.com/gin-gonic/gin"
)

func (w *WebRouter) SetupUserRouter(r *gin.RouterGroup) {
	// 用户管理相关API
	group := r.Group("/users")
	group.GET("", w.Handler.HandleListUsers)
	group.GET("/", w.Handler.HandleListUsers)              // 获取用户列表
	group.POST("/", w.Handler.HandleCreateUser)            // 创建用户
	group.DELETE("/:username", w.Handler.HandleDeleteUser) // 删除用户
}
