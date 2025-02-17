package routers

import (
	"servon/core/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupFileRouter(r *gin.RouterGroup) {
	api := r.Group("/files")
	api.GET("/", handlers.HandleFileList)
	api.GET("", handlers.HandleFileList)
	api.GET("/download", handlers.HandleFileDownload)
	api.GET("/content", handlers.HandleFileContent)
	api.POST("/save", handlers.HandleSaveFile)
	api.DELETE("/delete", handlers.HandleDeleteFile)
	api.POST("/create", handlers.HandleCreateFile)
	api.POST("/rename", handlers.HandleRenameFile)
}
