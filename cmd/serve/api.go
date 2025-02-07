package web

import (
	"servon/cmd/serve/handler"
)

// setupAPIRoutes 设置所有API路由
func (s *Server) setupAPIRoutes() {
	h := handler.New()
	sh := handler.NewSoftwareHandler()
	dh := handler.NewDeployHandler()
	api := s.router.Group("/web_api")
	{
		api.GET("/system/resources", h.HandleSystemResources)
		api.GET("/system/network", h.HandleNetworkResources)
		api.GET("/system/user", h.HandleCurrentUser)
		api.GET("/system/os", h.HandleOSInfo)
		api.GET("/system/basic", h.HandleBasicInfo)
		api.GET("/system/software", sh.HandleGetSoftwareList)
		api.GET("/system/software/:name/install", sh.HandleInstallSoftware)
		api.GET("/system/software/:name/uninstall", sh.HandleUninstallSoftware)
		api.POST("/system/software/:name/stop", sh.HandleStopSoftware)
		api.GET("/system/software/:name/status", sh.HandleGetSoftwareStatus)
		api.GET("/system/processes", h.HandleProcessList)
		api.GET("/system/files", h.HandleFileList)
		api.GET("/system/ports", h.HandlePortList)

		// 定时任务相关API
		api.GET("/cron/tasks", h.HandleListCronTasks)              // 获取所有定时任务
		api.POST("/cron/tasks", h.HandleCreateCronTask)            // 创建定时任务
		api.PUT("/cron/tasks/:id", h.HandleUpdateCronTask)         // 更新定时任务
		api.DELETE("/cron/tasks/:id", h.HandleDeleteCronTask)      // 删除定时任务
		api.POST("/cron/tasks/:id/toggle", h.HandleToggleCronTask) // 启用/禁用定时任务

		// 部署相关API
		api.GET("/deploy/projects", dh.HandleListProjects)            // 获取所有项目
		api.POST("/deploy/projects", dh.HandleCreateProject)          // 创建项目
		api.PUT("/deploy/projects/:id", dh.HandleUpdateProject)       // 更新项目
		api.DELETE("/deploy/projects/:id", dh.HandleDeleteProject)    // 删除项目
		api.POST("/deploy/projects/:id/build", dh.HandleBuildProject) // 构建项目
		api.GET("/deploy/projects/:id/logs", dh.HandleProjectLogs)    // 获取项目日志
	}
}
