package web

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"path"

	"github.com/gin-gonic/gin"
)

//go:embed dist
var distFS embed.FS

// setupUIRoutes 设置所有UI相关路由
func (s *Server) setupUIRoutes() {
	// 获取嵌入的dist子文件系统
	subFS, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}

	// 创建静态文件处理器
	// fileServer := http.FileServer(http.FS(subFS))

	// 处理静态资源请求
	s.router.GET("/assets/*filepath", func(c *gin.Context) {
		c.FileFromFS(path.Join("assets", c.Param("filepath")), http.FS(subFS))
	})
	fmt.Println("Registered static assets route at /assets/*filepath")

	// 在setupUIRoutes开始处添加：
	fmt.Println("===== Embedded Files =====")
	fs.WalkDir(subFS, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			fmt.Printf("|-- %s\n", path)
		}
		return nil
	})
	fmt.Println("==========================")

	// 2. 确保先注册静态路由再注册通配路由
	s.router.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		content, _ := fs.ReadFile(subFS, "index.html")
		c.Data(200, "text/html; charset=utf-8", content)
	})

	// 添加明确的静态文件路由
	s.router.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("/favicon.ico", http.FS(subFS))
	})

	// 3. 最后处理其他前端路由
	s.router.NoRoute(func(c *gin.Context) {
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
		content, _ := fs.ReadFile(subFS, "index.html")
		c.Data(200, "text/html; charset=utf-8", content)
	})
}
