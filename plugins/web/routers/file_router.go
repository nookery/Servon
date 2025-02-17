package routers

import (
	"github.com/gin-gonic/gin"
)

func (w *WebRouter) SetupFileRouter(r *gin.RouterGroup) {
	api := r.Group("/files")
	api.GET("/", w.Handler.HandleFileList)
	api.GET("", w.Handler.HandleFileList)
	api.GET("/download", w.Handler.HandleFileDownload)
	api.GET("/content", w.Handler.HandleFileContent)
	api.POST("/save", w.Handler.HandleSaveFile)
	api.DELETE("/delete", w.Handler.HandleDeleteFile)
	api.POST("/create", w.Handler.HandleCreateFile)
	api.POST("/rename", w.Handler.HandleRenameFile)
}
