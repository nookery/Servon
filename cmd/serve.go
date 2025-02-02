package cmd

import (
	"servon/internal/web"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	port    int
	apiOnly bool
)

// ServeCmd 表示 serve 命令
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动服务器",
	Long: `启动服务器，提供系统管理API。

可选参数：
  --port, -p:    指定服务器监听端口，默认为8080
  --api-only:   仅启动API服务，不提供Web界面

示例：
  servon serve              # 启动完整服务（API + Web界面）
  servon serve --api-only   # 仅启动API服务
  servon serve -p 3000      # 在3000端口启动完整服务`,

	RunE: func(cmd *cobra.Command, args []string) error {
		color.Cyan("Starting Servon server on port %d...\n", port)
		if !apiOnly {
			color.Cyan("Web UI enabled\n")
		} else {
			color.Yellow("Running in API-only mode\n")
		}

		server := web.NewServer(port, !apiOnly)
		return server.Start()
	},
}

func init() {
	// 配置 serve 子命令的参数
	ServeCmd.Flags().IntVarP(&port, "port", "p", 8080, "服务器监听端口")
	ServeCmd.Flags().BoolVar(&apiOnly, "api-only", false, "仅启动API服务")
}
