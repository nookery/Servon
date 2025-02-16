package routers

import (
	"servon/core/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupCronRouter(r *gin.RouterGroup) {
	// 定时任务相关API
	group := r.Group("/cron")
	group.GET("/tasks", handlers.HandleListCronTasks)              // 获取所有定时任务
	group.POST("/tasks", handlers.HandleCreateCronTask)            // 创建定时任务
	group.PUT("/tasks/:id", handlers.HandleUpdateCronTask)         // 更新定时任务
	group.DELETE("/tasks/:id", handlers.HandleDeleteCronTask)      // 删除定时任务
	group.POST("/tasks/:id/toggle", handlers.HandleToggleCronTask) // 启用/禁用定时任务
}
