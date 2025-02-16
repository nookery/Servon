package serve

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

// setupAPIRoutes 设置所有API路由
func (p *ServePlugin) setupAPIRoutes(router *gin.Engine) {
	p.PrintKeyValue("API:", color.HiGreenString("http://localhost:%d/web_api", port)) // 仅当监听非本地地址时显示网络访问信息
	if host != "127.0.0.1" && host != "localhost" {
		p.PrintKeyValue("Network:", color.HiGreenString("http://%s:%d", host, port))
	}
}

//go:embed dist
var distFS embed.FS

// setupUIRoutes 设置所有UI相关路由
func setupUIRoutes(router *gin.Engine) {
	// 获取嵌入的dist子文件系统
	subFS, err := fs.Sub(distFS, "dist")
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
