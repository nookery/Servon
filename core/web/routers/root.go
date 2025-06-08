package routers

import (
	"servon/core/managers"
	"servon/core/web/controllers"

	"github.com/gin-gonic/gin"
)

func Setup(manager *managers.FullManager, r *gin.Engine, isDev bool) {
	serviceController := controllers.NewServiceController(manager.ServiceManager)
	deployController := controllers.NewDeployController(manager)
	fileController := controllers.NewFileController(manager)
	cronController := controllers.NewCronController(manager)

	api := r.Group("/web_api")

	SetupSoftRouter(api, manager)
	SetupProcessRouter(api, manager)
	SetupInfoRouter(api, manager)
	SetupGitHubRouter(api, manager)
	SetupTaskRouter(api, manager)
	SetupPortRouter(api, manager)
	SetupUserRouter(api, manager)
	SetupIntegrationRouter(api, manager)
	SetupLogRouter(api, manager.LogManager)
	SetupTopologyRoutes(api, manager.ProjectManager)

	// 定时任务相关API
	group := r.Group("/cron")
	group.GET("/tasks", cronController.HandleListCronTasks)              // 获取所有定时任务
	group.POST("/tasks", cronController.HandleCreateCronTask)            // 创建定时任务
	group.PUT("/tasks/:id", cronController.HandleUpdateCronTask)         // 更新定时任务
	group.DELETE("/tasks/:id", cronController.HandleDeleteCronTask)      // 删除定时任务
	group.POST("/tasks/:id/toggle", cronController.HandleToggleCronTask) // 启用/禁用定时任务

	// 文件管理
	fileGroup := api.Group("/files")
	fileGroup.GET("/", fileController.HandleFileList)
	fileGroup.GET("", fileController.HandleFileList)
	fileGroup.GET("/download", fileController.HandleFileDownload)
	fileGroup.GET("/content", fileController.HandleFileContent)
	fileGroup.POST("/save", fileController.HandleSaveFile)
	fileGroup.DELETE("/delete", fileController.HandleDeleteFile)
	fileGroup.POST("/create", fileController.HandleCreateFile)
	fileGroup.POST("/rename", fileController.HandleRenameFile)
	fileGroup.POST("/batch-delete", fileController.HandleBatchDeleteFiles)
	fileGroup.POST("/copy", fileController.HandleCopyFile)

	// 部署管理
	deployRouter := api.Group("/deploy")
	deployRouter.POST("/repository", deployController.DeployRepository)

	// 服务管理路由组
	serviceGroup := api.Group("/services")
	{
		// 获取所有服务列表
		serviceGroup.GET("", serviceController.GetServiceList)

		// 服务操作
		serviceGroup.POST("/:name/start", serviceController.StartService)
		serviceGroup.POST("/:name/stop", serviceController.StopService)
		serviceGroup.POST("/:name/restart", serviceController.RestartService)

		// 服务配置
		serviceGroup.GET("/:name/config", serviceController.GetServiceConfig)
		serviceGroup.PUT("/:name/config", serviceController.UpdateServiceConfig)

		// 服务日志
		serviceGroup.GET("/:name/logs", serviceController.GetServiceLogs)

		// 服务详情
		serviceGroup.GET("/:name/details", serviceController.GetServiceDetails)

		// 添加后台服务
		serviceGroup.POST("", serviceController.AddBackgroundService)

		// 删除后台服务
		serviceGroup.DELETE("/:name", serviceController.RemoveBackgroundService)
	}

	SetupUIRoutes(r)
}
