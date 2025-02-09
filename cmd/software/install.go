package software

import (
	"fmt"
	"servon/utils/logger"

	"github.com/spf13/cobra"
)

// InstallCmd represents the install command
var InstallCmd = &cobra.Command{
	Use:   "install [软件名称]",
	Short: "安装指定的软件",
	Long: `安装指定的软件。

支持的软件:
  - caddy: Web服务器
  - node: Node.js运行时
  - pnpm: 快速的包管理器
  - npm: Node.js包管理器
  - clash: 代理工具`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("请指定要安装的软件名称")
		}

		name := args[0]
		manager := NewSoftwareManager()

		// 检查软件是否支持
		if !manager.IsSupportedSoftware(name) {
			return fmt.Errorf("不支持的软件: %s", name)
		}

		// 开始安装
		logger.InfoTitle("📦 开始安装 %s ...", name)

		err := manager.InstallSoftware(name, nil)
		if err != nil {
			return fmt.Errorf("安装失败: %v", err)
		}

		logger.InfoTitle("✅ %s 安装完成!", name)
		return nil
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return []string{"caddy", "node", "pnpm", "npm", "clash"}, cobra.ShellCompDirectiveNoFileComp
	},
}
