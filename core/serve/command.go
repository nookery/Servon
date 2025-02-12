package serve

import (
	"fmt"
	"os"
	"os/exec"
	"servon/core/libs"
	"strconv"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	port    int
	host    string
	apiOnly bool
)

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "启动服务器",
		Run: func(cmd *cobra.Command, args []string) {
			apiOnly, _ := cmd.Flags().GetBool("api-only")
			appEnv := libs.DefaultEnvManager.GetEnv("APP_ENV")

			// 清晰的启动横幅
			fmt.Printf("\n  %s\n\n", color.HiCyanString("SERVON"))

			printKeyValue("Version:", libs.DefaultVersionManager.GetVersion())
			printKeyValue("API Only:", color.HiGreenString("%t", apiOnly))
			printKeyValue("Environment:", color.HiGreenString("%s", appEnv))
			fmt.Println()

			StartWebServer(host, port, !apiOnly)
		},
	}

	defaultPort, err := strconv.Atoi(libs.DefaultEnvManager.GetEnv("WEB_PORT_DEFAULT"))
	if err != nil {
		libs.PrintErrorf("Error: 无法将 WEB_PORT_DEFAULT 转换为整数，WEB_PORT_DEFAULT 的值为 %s", libs.DefaultEnvManager.GetEnv("WEB_PORT_DEFAULT"))
		os.Exit(1)
	}

	// 配置 serve 子命令的参数
	cmd.Flags().IntVarP(&port, "port", "p", defaultPort, "服务器监听端口")
	cmd.Flags().StringVar(&host, "host", "127.0.0.1", "服务器监听地址")
	cmd.Flags().BoolVar(&apiOnly, "api-only", false, "仅启动API服务")

	return cmd
}

func printKeyValue(key string, value string) {
	fmt.Printf("  %-25s    %s\n",
		color.HiBlackString(key),
		color.HiGreenString(value))
}

func StartWebServer(host string, port int, withUI bool) {
	appEnv := libs.DefaultEnvManager.GetEnv("APP_ENV")
	router := libs.NewWebServer(host, port, withUI)

	// 设置API路由
	setupAPIRoutes(router)

	// 如果启用了UI，设置UI路由
	if withUI {
		// 检查是否为开发环境
		if appEnv == "development" {
			fmt.Println()
			libs.PrintInfof("开发环境，启动 npm dev server")
			libs.PrintInfof("VITE_API_TARGET=http://localhost:%d", port)
			go func() {
				cmd := exec.Command("npm", "run", "dev")
				cmd.Dir = "."
				cmd.Env = append(os.Environ(), "VITE_API_TARGET=http://localhost:"+strconv.Itoa(port))
				// 设置输出管道
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					fmt.Printf("Failed to start npm dev server: %v\n", err)
				}
			}()
		} else {
			setupUIRoutes(router)
			printKeyValue("Local UI:", color.HiGreenString("http://localhost:%d", port))
		}
	}

	router.Run(fmt.Sprintf("%s:%d", host, port))
}

// setupAPIRoutes 设置所有API路由
func setupAPIRoutes(router *gin.Engine) {
	h := NewHanlder()
	api := router.Group("/web_api")
	{
		api.GET("/system/resources", h.HandleSystemResources)
		api.GET("/system/network", h.HandleNetworkResources)
		api.GET("/system/user", h.HandleCurrentUser)
		api.GET("/system/os", h.HandleOSInfo)
		api.GET("/system/basic", h.HandleBasicInfo)
		api.GET("/system/software", h.HandleGetSoftwareList)
		api.GET("/system/software/:name/install", h.HandleInstallSoftware)
		api.GET("/system/software/:name/uninstall", h.HandleUninstallSoftware)
		api.POST("/system/software/:name/stop", h.HandleStopSoftware)
		api.GET("/system/software/:name/status", h.HandleGetSoftwareStatus)
		api.GET("/system/processes", h.HandleProcessList)
		api.GET("/system/files", h.HandleFileList)
		api.GET("/system/ports", h.HandlePortList)

		// // 定时任务相关API
		// api.GET("/cron/tasks", h.HandleListCronTasks)              // 获取所有定时任务
		// api.POST("/cron/tasks", h.HandleCreateCronTask)            // 创建定时任务
		// api.PUT("/cron/tasks/:id", h.HandleUpdateCronTask)         // 更新定时任务
		// api.DELETE("/cron/tasks/:id", h.HandleDeleteCronTask)      // 删除定时任务
		// api.POST("/cron/tasks/:id/toggle", h.HandleToggleCronTask) // 启用/禁用定时任务
	}

	printKeyValue("API:", color.HiGreenString("http://localhost:%d/web_api", port)) // 仅当监听非本地地址时显示网络访问信息
	if host != "127.0.0.1" && host != "localhost" {
		printKeyValue("Network:", color.HiGreenString("http://%s:%d", host, port))
	}
}
