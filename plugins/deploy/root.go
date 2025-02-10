package deploy

import (
	"github.com/spf13/cobra"
	"servon/core"
)

func Setup(core *core.Core) {
	core.AddCommand(GetDeployCommand(core))
}

// GetDeployCommand 返回 deploy 命令
func GetDeployCommand(core *core.Core) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "部署项目",
		Long: `部署项目。

示例：
  servon deploy start    # 启动部署
  servon deploy stop     # 停止部署`,
		RunE: func(cmd *cobra.Command, args []string) error {
			core.PrintCommandHelp(cmd)
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
