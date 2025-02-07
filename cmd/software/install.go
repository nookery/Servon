package software

import (
	"fmt"
	"servon/cmd/utils/logger"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// newInstallCmd 返回 install 子命令
func newInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install [软件名称]",
		Short: "安装指定的软件",
		Long: `安装指定的软件到系统中。

示例：
  servon software install nginx    # 安装 nginx
  servon software install mysql    # 安装 mysql`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				color.New(color.FgRed).Println("\n❌ 缺少软件名称参数")
				fmt.Println("\n用法:")
				color.New(color.FgYellow).Print("  servon software install ")
				fmt.Println("[软件名称]")

				// 显示支持的软件列表
				manager := NewSoftwareManager()
				names := manager.GetSoftwareNames()
				fmt.Println("\n支持的软件:")
				for _, name := range names {
					color.New(color.FgHiWhite).Printf("  - %s\n", name)
				}

				fmt.Println("\n示例:")
				color.New(color.FgCyan).Println("  servon software install nginx")
				color.New(color.FgCyan).Println("  servon software install mysql")
				return nil
			}

			manager := NewSoftwareManager()
			name := args[0]

			// 检查软件是否支持
			supported := false
			for _, sw := range manager.GetSoftwareNames() {
				if sw == name {
					supported = true
					break
				}
			}

			if !supported {
				color.New(color.FgRed).Printf("\n❌ 不支持的软件: %s\n", name)
				fmt.Println("\n支持的软件:")
				for _, sw := range manager.GetSoftwareNames() {
					color.New(color.FgHiWhite).Printf("  - %s\n", sw)
				}
				return nil
			}

			// 开始安装
			logger.InfoTitle("📦 开始安装 %s ...", name)

			err := manager.InstallSoftware(name, nil)
			if err != nil {
				color.New(color.FgRed).Printf("\n❌ 安装失败: %v\n", err)
				return nil
			}

			logger.InfoTitle("✨ 软件 %s 安装完成！", name)

			return nil
		},
	}
}
