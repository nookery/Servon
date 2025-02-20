package commands

import (
	"os"
	"os/exec"
	"servon/core/internal/libs/managers"
	"servon/core/internal/libs/utils"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
	port    int
	host    string
	apiOnly bool
	devMode bool
)

// GetServerCommand 获取服务器管理命令
func GetServerCommand(web *utils.WebServer, manager *managers.FullManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "启动服务器",
	}

	cmd.AddCommand(MakeStartCommand(web, manager))
	cmd.AddCommand(MakeStopCommand(web))
	cmd.AddCommand(MakeRestartCommand(web))

	return cmd
}

// MakeStartCommand 创建启动命令
func MakeStartCommand(web *utils.WebServer, manager *managers.FullManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "启动服务器",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetInt("port")
			host, _ := cmd.Flags().GetString("host")
			apiOnly, _ := cmd.Flags().GetBool("api-only")
			devMode, _ := cmd.Flags().GetBool("dev")

			web.SetPort(port)
			web.SetHost(host)

			if devMode {
				printer.PrintInfof("开发环境，先关闭服务器")
				if err := web.StopBackground(); err != nil {
					printer.PrintError(err)
					os.Exit(1)
				}
			}

			// 使用 RunUntilSignal 来保持服务器运行
			if err := web.RunInBackground(); err != nil {
				printer.PrintError(err)
				os.Exit(1)
			}

			printer.PrintSuccess("服务器启动成功")

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
				"Dev Mode":  stringUtil.GetEmojiForBool(devMode),
				"Link":      "http://" + web.GetHost() + ":" + web.GetPortString(),
				"Notice":    "服务器在后台运行，如需要关闭，执行: servon serve stop",
			})
			printer.PrintLn()

			if devMode {
				printer.PrintInfof("开发环境，启动 npm dev server")
				runNpmDev(web.GetPort())
			}
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 9754, "服务器监听端口")
	cmd.Flags().StringVar(&host, "host", "127.0.0.1", "服务器监听地址")
	cmd.Flags().BoolVar(&apiOnly, "api-only", false, "仅启动API服务")
	cmd.Flags().BoolVar(&devMode, "dev", false, "启用开发模式")

	return cmd
}

// MakeStopCommand 创建停止命令
func MakeStopCommand(web *utils.WebServer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "停止服务器",
		Run: func(cmd *cobra.Command, args []string) {
			if err := web.StopBackground(); err != nil {
				printer.PrintError(err)
				os.Exit(1)
			}

		},
	}

	return cmd
}

// MakeRestartCommand 创建重启命令
func MakeRestartCommand(web *utils.WebServer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "重启服务器",
		Run: func(cmd *cobra.Command, args []string) {
			// 获取配置
			port, _ := cmd.Flags().GetInt("port")
			host, _ := cmd.Flags().GetString("host")

			// 设置配置
			web.SetPort(port)
			web.SetHost(host)

			// 先停止
			if err := web.StopBackground(); err != nil {
				printer.PrintWarningf("停止服务器时出错: %v", err)
			}

			// 等待一小段时间确保端口释放
			time.Sleep(time.Second)

			// 重新启动
			if err := web.RunInBackground(); err != nil {
				printer.PrintError(err)
				os.Exit(1)
			}

			printer.PrintSuccess("服务器已重启")
		},
	}

	// 添加与 start 命令相同的标志
	cmd.Flags().IntVarP(&port, "port", "p", 9754, "服务器监听端口")
	cmd.Flags().StringVar(&host, "host", "127.0.0.1", "服务器监听地址")

	return cmd
}

// MakeDevCommand 创建开发命令
func MakeDevCommand(web *utils.WebServer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "启动开发服务器",
		Run: func(cmd *cobra.Command, args []string) {
			// 获取配置
			port, _ := cmd.Flags().GetInt("port")
			host, _ := cmd.Flags().GetString("host")

			// 设置配置
			web.SetPort(port)
			web.SetHost(host)

			// 先停止
			if err := web.StopBackground(); err != nil {
				printer.PrintWarningf("停止服务器时出错: %v", err)
			}

			// 等待一小段时间确保端口释放
			time.Sleep(time.Second)

			// 重新启动
			if err := web.RunInBackground(); err != nil {
				printer.PrintError(err)
				os.Exit(1)
			}

			printer.PrintSuccess("服务器已重启")
		},
	}

	// 添加与 start 命令相同的标志
	cmd.Flags().IntVarP(&port, "port", "p", 9754, "服务器监听端口")
	cmd.Flags().StringVar(&host, "host", "127.0.0.1", "服务器监听地址")

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
