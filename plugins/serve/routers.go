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
	api := router.Group("/web_api")
	{
		api.GET("/system/resources", p.HandleSystemResources)
		api.GET("/system/network", p.HandleNetworkResources)
		api.GET("/system/user", p.HandleCurrentUser)
		api.GET("/system/os", p.HandleOSInfo)
		api.GET("/system/basic", p.HandleBasicInfo)

		// 软件管理相关API
		api.GET("/system/software", p.HandleGetSoftwareList)
		api.POST("/system/software/:name/install", p.HandleInstallSoftware)
		api.POST("/system/software/:name/uninstall", p.HandleUninstallSoftware)
		api.POST("/system/software/:name/stop", p.HandleStopSoftware)
		api.POST("/system/software/:name/start", p.HandleStartSoftware)
		api.GET("/system/software/:name/status", p.HandleGetSoftwareStatus)

		// 进程管理相关API
		api.GET("/system/processes", p.HandleProcessList)

		// 文件管理相关API
		api.GET("/system/files", p.HandleFileList)
		api.GET("/system/files/download", p.HandleFileDownload)
		api.GET("/system/files/content", p.HandleFileContent)
		api.POST("/system/files/save", p.HandleSaveFile)
		api.DELETE("/system/files/delete", p.HandleDeleteFile)
		api.POST("/system/files/create", p.HandleCreateFile)

		// 端口管理相关API
		api.GET("/system/ports", p.HandlePortList)

		// Add new streaming logs endpoint
		api.GET("/logs/:channel", p.HandleStreamLogs)

		// 定时任务相关API
		api.GET("/cron/tasks", p.HandleListCronTasks)              // 获取所有定时任务
		api.POST("/cron/tasks", p.HandleCreateCronTask)            // 创建定时任务
		api.PUT("/cron/tasks/:id", p.HandleUpdateCronTask)         // 更新定时任务
		api.DELETE("/cron/tasks/:id", p.HandleDeleteCronTask)      // 删除定时任务
		api.POST("/cron/tasks/:id/toggle", p.HandleToggleCronTask) // 启用/禁用定时任务

		// 用户管理相关API
		api.GET("/users", p.HandleListUsers)               // 获取用户列表
		api.POST("/users", p.HandleCreateUser)             // 创建用户
		api.DELETE("/users/:username", p.HandleDeleteUser) // 删除用户
	}

	printKeyValue("API:", color.HiGreenString("http://localhost:%d/web_api", port)) // 仅当监听非本地地址时显示网络访问信息
	if host != "127.0.0.1" && host != "localhost" {
		printKeyValue("Network:", color.HiGreenString("http://%s:%d", host, port))
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
