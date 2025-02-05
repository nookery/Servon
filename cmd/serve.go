package cmd

import (
	"fmt"
	web "servon/internal/serve"
	"servon/internal/version"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	port    int
	apiOnly bool
	host    string
)

// ServeCmd 表示 serve 命令
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动服务器",
	Long: `启动服务器，提供系统管理API。

可选参数：
  --port, -p:     指定服务器监听端口，默认为8080
  --host:         指定服务器监听地址，默认为127.0.0.1，使用0.0.0.0监听所有地址
  --api-only:     仅启动API服务，不提供Web界面

示例：
  servon serve                    # 在127.0.0.1:8080启动完整服务
  servon serve --host 0.0.0.0    # 监听所有地址
  servon serve --api-only        # 仅启动API服务
  servon serve -p 3000           # 在3000端口启动完整服务`,

	RunE: func(cmd *cobra.Command, args []string) error {
		// 清晰的启动横幅
		fmt.Printf("\n  %s\n\n", color.HiCyanString("SERVON"))

		// 版本和模式信息
		fmt.Printf("  %s    %s\n",
			color.HiBlackString("Version:"),
			color.HiWhiteString(version.GetVersion()))
		fmt.Printf("  %s    %s\n",
			color.HiBlackString("Mode:"),
			color.HiWhiteString(map[bool]string{true: "API Only", false: "Full Stack"}[apiOnly]))

		// 访问信息
		fmt.Printf("  %s    %s\n",
			color.HiBlackString("Local:"),
			color.HiGreenString("http://localhost:%d", port))

		// 仅当监听非本地地址时显示网络访问信息
		if host != "127.0.0.1" && host != "localhost" {
			fmt.Printf("  %s    %s\n",
				color.HiBlackString("Network:"),
				color.HiGreenString("http://%s:%d", host, port))
		}
		fmt.Printf("\n")

		if !apiOnly {
			fmt.Printf("  %s\n", color.HiBlackString("Web UI:"))
			fmt.Printf("    • Dashboard    %s\n", color.HiGreenString("http://localhost:%d", port))
			fmt.Printf("    • API Docs     %s\n", color.HiGreenString("http://localhost:%d/docs", port))
		}

		fmt.Printf("\n  %s  %s\n\n",
			color.YellowString("⚡"),
			color.HiBlackString("Server is ready"))

		server := web.NewServer(host, port, !apiOnly)
		return server.Start()
	},
}

func init() {
	// 配置 serve 子命令的参数
	ServeCmd.Flags().IntVarP(&port, "port", "p", 8080, "服务器监听端口")
	ServeCmd.Flags().StringVar(&host, "host", "127.0.0.1", "服务器监听地址")
	ServeCmd.Flags().BoolVar(&apiOnly, "api-only", false, "仅启动API服务")
}
