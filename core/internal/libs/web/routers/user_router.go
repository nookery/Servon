package routers

import (
	"servon/core/internal/libs/managers"
	"servon/core/internal/libs/web/controllers"

	"github.com/gin-gonic/gin"
)

func SetupUserRouter(r *gin.RouterGroup, manager *managers.FullManager) {
	controller := controllers.NewUserController(manager)

	// 用户管理相关API
	group := r.Group("/users")
	group.GET("", controller.HandleListUsers)
	group.GET("/", controller.HandleListUsers)              // 获取用户列表
	group.POST("/", controller.HandleCreateUser)            // 创建用户
	group.DELETE("/:username", controller.HandleDeleteUser) // 删除用户
}
