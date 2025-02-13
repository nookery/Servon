package serve

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
func (p *ServePlugin) RegisterRoutes(router *gin.RouterGroup) {
	// ... 其他已有路由 ...

	// 文件管理相关路由
	router.GET("/files", p.HandleFileList)
	router.GET("/files/download", p.HandleFileDownload)
	router.GET("/files/content", p.HandleFileContent)
	router.POST("/files/save", p.HandleSaveFile)
	router.DELETE("/files/delete", p.HandleDeleteFile)
	router.POST("/files/create", p.HandleCreateFile)
}
