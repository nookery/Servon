package routers

import (
	"servon/core/managers"
	"servon/plugins/server/web/controllers"

	"github.com/gin-gonic/gin"
)

func SetupLogRouter(r *gin.RouterGroup, logManager *managers.LogManager) {
	controller := controllers.NewLogController(logManager)

	group := r.Group("/logs")
	{
		group.GET("/files", controller.HandleListLogFiles)     // 获取日志文件列表
		group.GET("/entries", controller.HandleReadLogEntries) // 读取日志内容
		group.GET("/search", controller.HandleSearchLogs)      // 搜索日志
		group.GET("/stats", controller.HandleGetLogStats)      // 获取日志统计
		group.POST("/clean", controller.HandleCleanOldLogs)    // 清理旧日志
		group.POST("/delete", controller.HandleDeleteLogFile)  // 删除指定日志文件
		group.POST("/clear", controller.HandleClearLogFile)    // 清空指定日志文件
	}
}
