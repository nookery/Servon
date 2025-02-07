package deploy

import (
	"servon/cmd/internal/utils"

	"github.com/spf13/cobra"
)

// GetDeployCommand 返回 deploy 命令
func GetDeployCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "部署项目",
		Long: `部署项目。

示例：
  servon deploy start    # 启动部署
  servon deploy stop     # 停止部署`,
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.PrintCommandHelp(cmd, map[string]string{
				"start": "启动部署",
				"stop":  "停止部署",
				"get":   "获取项目列表",
				"serve": "创建静态文件服务",
			})
			return nil
		},
	}

	// 添加子命令
	cmd.AddCommand(newStartCmd())
	cmd.AddCommand(newStopCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newServeCmd())

	// 设置 serve 命令的参数
	serveCmd := cmd.Commands()[3]
	serveCmd.Flags().String("name", "", "服务名称")
	serveCmd.Flags().String("path", "", "本地文件夹路径")
	serveCmd.Flags().String("domain", "", "访问域名")

	return cmd
}
