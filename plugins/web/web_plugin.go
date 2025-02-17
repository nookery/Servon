package web

import (
	"servon/core"
	"servon/plugins/web/routers"

	"github.com/spf13/cobra"
)

type webPlugin struct {
	*core.App
	githubStates []string // 用于存储GitHub integration的state值
}

func Setup(app *core.App) {
	plugin := &webPlugin{
		App:          app,
		githubStates: make([]string, 0),
	}

	app.GetRootCommand().AddCommand(plugin.newServeCommand())
}

func (p *webPlugin) newServeCommand() *cobra.Command {
	var (
		port    int
		host    string
		apiOnly bool
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "启动服务器",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetInt("port")
			host, _ := cmd.Flags().GetString("host")
			apiOnly, _ := cmd.Flags().GetBool("api-only")

			p.App.WebServerManager.SetPort(port)
			p.App.WebServerManager.SetHost(host)

			routers.Setup(p.App, p.IsDev())

			// 启动横幅
			p.PrintLn()
			p.PrintTitle("SERVON")
			p.PrintLn()
			p.PrintKeyValues(map[string]string{
				"Version":   p.App.VersionManager.GetVersion(),
				"API Route": p.App.StringUtil.GetEmojiForBool(true),
				"UI Route":  p.App.StringUtil.GetEmojiForBool(!apiOnly),
				"Port":      p.GetPortString(),
				"Host":      p.GetHost(),
			})
			p.PrintLn()

			p.StartWebServer()
		},
	}

	// 配置 serve 子命令的参数
	cmd.Flags().IntVarP(&port, "port", "p", 9754, "服务器监听端口")
	cmd.Flags().StringVar(&host, "host", "127.0.0.1", "服务器监听地址")
	cmd.Flags().BoolVar(&apiOnly, "api-only", false, "仅启动API服务")

	return cmd
}
