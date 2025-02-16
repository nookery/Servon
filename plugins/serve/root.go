package serve

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"servon/core"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	port    int
	host    string
	apiOnly bool
)

// Config represents the configuration for the serve plugin
type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type ServePlugin struct {
	*core.App
	Config       *Config
	githubStates []string // 用于存储GitHub integration的state值
}

func Setup(app *core.App) {
	plugin := &ServePlugin{
		App: app,
		Config: &Config{
			Host: "43.142.208.212",
			Port: 9754,
		},
		githubStates: make([]string, 0),
	}
	app.GetRootCommand().AddCommand(plugin.NewServeCommand())
}

func (p *ServePlugin) NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "启动服务器",
		Run: func(cmd *cobra.Command, args []string) {
			apiOnly, _ := cmd.Flags().GetBool("api-only")

			// 清晰的启动横幅
			fmt.Printf("\n  %s\n\n", color.HiCyanString("SERVON"))

			p.PrintKeyValue("Version:", p.VersionManager.GetVersion())
			p.PrintKeyValue("API Only:", color.HiGreenString("%t", apiOnly))
			fmt.Println()

			p.StartWebServer(host, port, !apiOnly)
		},
	}

	// 配置 serve 子命令的参数
	cmd.Flags().IntVarP(&port, "port", "p", 9754, "服务器监听端口")
	cmd.Flags().StringVar(&host, "host", "127.0.0.1", "服务器监听地址")
	cmd.Flags().BoolVar(&apiOnly, "api-only", false, "仅启动API服务")

	return cmd
}

func (p *ServePlugin) StartWebServer(host string, port int, withUI bool) {
	router := p.GetRouter()

	// 设置API路由
	p.setupAPIRoutes(router)

	// 如果启用了UI，设置UI路由
	if withUI {
		// 检查是否为开发环境（通过检查是否使用 go run 启动）
		if os.Args[0] == "main" || strings.Contains(os.Args[0], "go-build") {
			go func() {
				p.PrintLn()
				p.PrintInfof("开发环境，启动 npm dev server")
				p.PrintInfof("VITE_API_TARGET=http://127.0.0.1:%d", port)
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
			p.PrintKeyValue("Local UI:", color.HiGreenString("http://localhost:%d", port))
		}
	}

	// 启动 Web 服务器
	// p.PrintInfof("启动 Web 服务器: http://%s:%d", host, port)
	err := router.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		p.PrintErrorf("启动 Web 服务器失败: %v", err)
		os.Exit(1)
	}
}
