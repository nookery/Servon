package routers

import (
	"servon/core/managers"
	"servon/plugins/server/web/controllers"

	"github.com/gin-gonic/gin"
)

func SetupTaskRouter(r *gin.RouterGroup, manager *managers.FullManager) {
	controller := controllers.NewTaskController(manager)

	// 任务管理相关API
	group := r.Group("/tasks")
	group.GET("", controller.HandleListTasks)
	group.GET("/", controller.HandleListTasks)               // 获取任务列表
	group.DELETE("/:id", controller.HandleRemoveTask)        // 删除任务
	group.POST("/:id/execute", controller.HandleExecuteTask) // 执行任务
}
