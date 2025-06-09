package server

import (
	"os"
	"os/exec"
	"servon/components/web_server"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// MakeDevCommand 创建开发命令
func MakeDevCommand(web *web_server.WebServer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "启动开发服务器",
		Run: func(cmd *cobra.Command, args []string) {
			// 先停止
			if err := web.StopBackground(); err != nil {
				log.Warnf("停止服务器时出错: %v", err)
			}

			// 等待一小段时间确保端口释放
			time.Sleep(time.Second)

			// 重新启动
			if err := web.RunInBackground(); err != nil {
				log.Error(err)
				os.Exit(1)
			}

			log.Success("服务器已重启")
		},
	}

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
		log.Errorf("启动 npm dev server 失败: %v", err)
	}
}
