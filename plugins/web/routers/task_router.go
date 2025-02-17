package routers

import (
	"github.com/gin-gonic/gin"
)

func (w *WebRouter) SetupTaskRouter(r *gin.RouterGroup) {
	// 任务管理相关API
	group := r.Group("/tasks")
	group.GET("", w.Handler.HandleListTasks)
	group.GET("/", w.Handler.HandleListTasks)               // 获取任务列表
	group.DELETE("/:id", w.Handler.HandleRemoveTask)        // 删除任务
	group.POST("/:id/execute", w.Handler.HandleExecuteTask) // 执行任务
}
