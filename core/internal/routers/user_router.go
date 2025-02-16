package routers

import (
	"servon/core/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupUserRouter(r *gin.RouterGroup) {
	// 用户管理相关API
	group := r.Group("/users")
	group.GET("", handlers.HandleListUsers)
	group.GET("/", handlers.HandleListUsers)              // 获取用户列表
	group.POST("/", handlers.HandleCreateUser)            // 创建用户
	group.DELETE("/:username", handlers.HandleDeleteUser) // 删除用户
}
