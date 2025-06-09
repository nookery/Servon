package server

import (
	"os"

	"github.com/spf13/cobra"
)

// MakeStartCommand 创建启动命令
func MakeStartCommand() *cobra.Command {
	var (
		port    int
		host    string
		apiOnly bool
		devMode bool
		verbose bool
	)
	cmd := &cobra.Command{
		Use:   "start",
		Short: "启动服务器",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ = cmd.Flags().GetInt("port")
			host, _ = cmd.Flags().GetString("host")
			apiOnly, _ = cmd.Flags().GetBool("api-only")
			devMode, _ = cmd.Flags().GetBool("dev")
			verbose, _ = cmd.Flags().GetBool("verbose")

			server.SetPort(port)
			server.SetHost(host)
			server.SetVerbose(verbose)

			if devMode {
				log.Infof("开发环境，先关闭服务器")
				if err := server.StopBackground(); err != nil {
					log.Error(err)
					os.Exit(1)
				}
			}

			// 使用 RunUntilSignal 来保持服务器运行
			if err := server.RunInBackground(); err != nil {
				log.Error(err)
				os.Exit(1)
			}

			log.Success("服务器启动成功")

			// 启动横幅
			log.EmptyLine()
			log.Title("SERVON")
			log.EmptyLine()
			log.PrintKeyValues(map[string]string{
				"API Route": stringUtil.GetEmojiForBool(true),
				"UI Route":  stringUtil.GetEmojiForBool(!apiOnly),
				"Port":      server.GetPortString(),
				"Host":      server.GetHost(),
				"Dev Mode":  stringUtil.GetEmojiForBool(devMode),
				"Link":      "http://" + server.GetHost() + ":" + server.GetPortString(),
				"Notice":    "服务器在后台运行，如需要关闭，执行: servon serve stop",
			})
			log.EmptyLine()

			if devMode {
				log.Infof("开发环境，启动 npm dev server")
				runNpmDev(server.GetPort())
			}
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 0, "服务器监听端口")
	cmd.Flags().StringVar(&host, "host", "127.0.0.1", "服务器监听地址")
	cmd.Flags().BoolVar(&apiOnly, "api-only", false, "仅启动API服务")
	cmd.Flags().BoolVar(&devMode, "dev", false, "启用开发模式")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "启用详细日志模式，显示端口检查和进程查找的详细信息")

	return cmd
}
