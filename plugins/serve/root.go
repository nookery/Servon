package serve

import (
	"fmt"
	"os"
	"os/exec"
	"servon/core"
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

type ServePlugin struct {
	*core.Core
}

func Setup(core *core.Core) {
	plugin := &ServePlugin{
		Core: core,
	}
	core.GetRootCommand().AddCommand(plugin.NewServeCommand())
}

func (p *ServePlugin) NewServeCommand() *cobra.Command {
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

			p.StartWebServer(host, port, !apiOnly)
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

func (p *ServePlugin) StartWebServer(host string, port int, withUI bool) {
	appEnv := libs.DefaultEnvManager.GetEnv("APP_ENV")
	router := libs.NewWebServer(host, port, withUI)

	// 设置API路由
	p.setupAPIRoutes(router)

	// 如果启用了UI，设置UI路由
	if withUI {
		// 检查是否为开发环境
		if appEnv == "development" {
			fmt.Println()
			libs.PrintInfof("开发环境，启动 npm dev server")
			libs.PrintInfof("VITE_API_TARGET=http://127.0.0.1:%d", port)
			go func() {
				cmd := exec.Command("npm", "run", "dev")
				cmd.Dir = "."
				cmd.Env = append(os.Environ(), "VITE_API_TARGET=http://127.0.0.1:"+strconv.Itoa(port))
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

// HandleStreamLogs streams logs from a specified channel using Server-Sent Events (SSE)
func (p *ServePlugin) HandleStreamLogs(c *gin.Context) {
	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// Create a channel to notify if client disconnects
	clientGone := c.Writer.CloseNotify()

	// Get the log channel (you'll need to implement this based on your logging system)
	logChan := p.LogChan
	if logChan == nil {
		c.String(404, "Log channel not found")
		return
	}

	// Stream logs
	for {
		select {
		case <-clientGone:
			// Client disconnected
			return
		case msg, ok := <-logChan:
			if !ok {
				// Channel closed
				return
			}
			// Send log message as SSE
			c.SSEvent("log", msg)
			c.Writer.Flush()
		}
	}
}
