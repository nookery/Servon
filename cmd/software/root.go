package software

import (
	"servon/cmd/utils"

	"github.com/spf13/cobra"
)

// GetSoftwareCommand 返回 software 命令
func GetSoftwareCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "software",
		Short: "软件管理",
		Long: `软件管理。

示例：
  servon software list      # 显示支持的软件列表
  servon software install   # 安装指定的软件
  servon software info      # 显示软件详细信息
  servon software start     # 启动指定的软件
  servon software stop      # 停止指定的软件`,
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.PrintCommandHelp(cmd, map[string]string{
				"list":    "显示支持的软件列表",
				"install": "安装指定的软件",
				"info":    "显示软件详细信息",
				"start":   "启动指定的软件",
				"stop":    "停止指定的软件",
			})
			return nil
		},
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(InstallCmd)
	cmd.AddCommand(newInfoCmd())
	cmd.AddCommand(newStartCmd())
	cmd.AddCommand(newStopCmd())

	return cmd
}
