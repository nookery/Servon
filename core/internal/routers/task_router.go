package routers

import (
	"servon/core/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupTaskRouter(r *gin.RouterGroup) {
	// 任务管理相关API
	group := r.Group("/tasks")
	group.GET("/", handlers.HandleListTasks)               // 获取任务列表
	group.DELETE("/:id", handlers.HandleRemoveTask)        // 删除任务
	group.POST("/:id/execute", handlers.HandleExecuteTask) // 执行任务
}
