package commands

import (
	"os"
	"os/exec"
	"servon/core/internal/libs/managers"
	"servon/core/internal/libs/utils"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	port    int
	host    string
	apiOnly bool
)

// GetUserRootCommand 获取用户管理命令
func GetServeCommand(web *utils.WebServer, manager *managers.FullManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "启动服务器",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetInt("port")
			host, _ := cmd.Flags().GetString("host")
			apiOnly, _ := cmd.Flags().GetBool("api-only")

			web.SetPort(port)
			web.SetHost(host)

			// 启动横幅
			printer.PrintLn()
			printer.PrintTitle("SERVON")
			printer.PrintLn()
			printer.PrintKeyValues(map[string]string{
				"Version":   manager.VersionManager.GetVersion(),
				"API Route": stringUtil.GetEmojiForBool(true),
				"UI Route":  stringUtil.GetEmojiForBool(!apiOnly),
				"Port":      web.GetPortString(),
				"Host":      web.GetHost(),
			})
			printer.PrintLn()

			printer.PrintInfof("正在启动服务器...")
			if err := web.Start(); err != nil {
				printer.PrintErrorf("%v", err)
				os.Exit(1)
			}

			if utils.DefaultDevUtil.IsDev() {
				printer.PrintInfof("开发环境，启动 npm dev server")
				runNpmDev(web.GetPort())
			}
		},
	}

	// 配置 serve 子命令的参数
	cmd.Flags().IntVarP(&port, "port", "p", 9754, "服务器监听端口")
	cmd.Flags().StringVar(&host, "host", "127.0.0.1", "服务器监听地址")
	cmd.Flags().BoolVar(&apiOnly, "api-only", false, "仅启动API服务")

	return cmd
}

func runNpmDev(backendPort int) {
	cmd := exec.Command("npm", "run", "dev")
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), "VITE_API_TARGET=http://127.0.0.1:"+strconv.Itoa(backendPort))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		printer.PrintErrorf("启动 npm dev server 失败: %v", err)
	}
}
