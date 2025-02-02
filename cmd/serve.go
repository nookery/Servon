package cmd

import (
	"servon/internal/web"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	port   int
	withUI bool
)

// RegisterCommands 注册所有子命令到根命令
func RegisterCommands(root *cobra.Command) {
	root.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动服务器",
	Long: `启动服务器，提供系统管理API。

可选参数：
  --port, -p: 指定服务器监听端口，默认为8080
  --ui, -u:   启用Web界面，从./dist目录提供前端文件

示例：
  servon serve              # 仅启动API服务
  servon serve --ui         # 启动API服务和Web界面
  servon serve -p 3000 -u   # 在3000端口启动完整服务`,

	RunE: func(cmd *cobra.Command, args []string) error {
		color.Cyan("Starting Servon server on port %d...\n", port)
		if withUI {
			color.Cyan("Web UI enabled\n")
		}

		server := web.NewServer(port, withUI)
		return server.Start()
	},
}

func init() {
	// 配置 serve 子命令的参数
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "服务器监听端口")
	serveCmd.Flags().BoolVarP(&withUI, "ui", "u", false, "是否提供Web界面")
}
