package routers

import (
	"github.com/gin-gonic/gin"
)

func (w *WebRouter) SetupCronRouter(r *gin.RouterGroup) {
	// 定时任务相关API
	group := r.Group("/cron")
	group.GET("/tasks", w.Handler.HandleListCronTasks)              // 获取所有定时任务
	group.POST("/tasks", w.Handler.HandleCreateCronTask)            // 创建定时任务
	group.PUT("/tasks/:id", w.Handler.HandleUpdateCronTask)         // 更新定时任务
	group.DELETE("/tasks/:id", w.Handler.HandleDeleteCronTask)      // 删除定时任务
	group.POST("/tasks/:id/toggle", w.Handler.HandleToggleCronTask) // 启用/禁用定时任务
}
