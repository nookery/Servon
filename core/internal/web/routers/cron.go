package routers

import (
	"servon/core/internal/managers"
	"servon/core/internal/web/controllers"

	"github.com/gin-gonic/gin"
)

func SetupCronRouter(r *gin.RouterGroup, manager *managers.FullManager) {
	cronController := controllers.NewCronController(manager)

	// 定时任务相关API
	group := r.Group("/cron")
	group.GET("/tasks", cronController.HandleListCronTasks)              // 获取所有定时任务
	group.POST("/tasks", cronController.HandleCreateCronTask)            // 创建定时任务
	group.PUT("/tasks/:id", cronController.HandleUpdateCronTask)         // 更新定时任务
	group.DELETE("/tasks/:id", cronController.HandleDeleteCronTask)      // 删除定时任务
	group.POST("/tasks/:id/toggle", cronController.HandleToggleCronTask) // 启用/禁用定时任务
}
