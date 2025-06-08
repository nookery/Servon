package routers

import (
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"servon/core/web/static"

	"github.com/gin-gonic/gin"
)

// SetupUIRoutes 设置所有UI相关路由
func SetupUIRoutes(router *gin.Engine) {
	// 获取嵌入的dist子文件系统
	subFS, err := fs.Sub(static.DistFS, "dist")
	if err != nil {
		panic(err)
	}

	// 处理静态资源请求
	router.GET("/assets/*filepath", func(c *gin.Context) {
		c.FileFromFS(path.Join("assets", c.Param("filepath")), http.FS(subFS))
	})

	// 2. 确保先注册静态路由再注册通配路由
	router.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		content, err := fs.ReadFile(subFS, "index.html")
		if err != nil {
			// 如果找不到 index.html，尝试读取 placeholder.html
			content, err = fs.ReadFile(subFS, "placeholder.html")
			if err != nil {
				c.String(http.StatusInternalServerError, "无法加载页面")
				return
			}
		}
		c.Data(200, "text/html; charset=utf-8", content)
	})

	// 添加明确的静态文件路由
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("/favicon.ico", http.FS(subFS))
	})

	// 3. 最后处理其他前端路由
	router.NoRoute(func(c *gin.Context) {
		fmt.Printf("Handling NoRoute: %s\n", c.Request.URL.Path)
		path := c.Request.URL.Path

		// 排除API和静态资源请求
		if strings.HasPrefix(path, "/web_api/") ||
			strings.HasPrefix(path, "/assets/") ||
			path == "/favicon.ico" {
			c.Next()
			return
		}

		c.Header("Content-Type", "text/html")
		content, err := fs.ReadFile(subFS, "index.html")
		if err != nil {
			// 如果找不到 index.html，尝试读取 placeholder.html
			content, err = fs.ReadFile(subFS, "placeholder.html")
			if err != nil {
				c.String(http.StatusInternalServerError, "无法加载页面")
				return
			}
		}
		c.Data(200, "text/html; charset=utf-8", content)
	})
}
