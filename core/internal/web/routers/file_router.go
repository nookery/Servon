package routers

import (
	"servon/core/internal/managers"
	"servon/core/internal/web/controllers"

	"github.com/gin-gonic/gin"
)

func SetupFileRouter(r *gin.RouterGroup, manager *managers.FullManager) {
	fileController := controllers.NewFileController(manager)

	api := r.Group("/files")
	api.GET("/", fileController.HandleFileList)
	api.GET("", fileController.HandleFileList)
	api.GET("/download", fileController.HandleFileDownload)
	api.GET("/content", fileController.HandleFileContent)
	api.POST("/save", fileController.HandleSaveFile)
	api.DELETE("/delete", fileController.HandleDeleteFile)
	api.POST("/create", fileController.HandleCreateFile)
	api.POST("/rename", fileController.HandleRenameFile)
	api.POST("/batch-delete", fileController.HandleBatchDeleteFiles)
}
