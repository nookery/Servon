package software

import (
	"fmt"
	"servon/internal/softwares"
	"servon/internal/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// newStartCmd 返回 start 子命令
func newStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start [软件名称]",
		Short: "启动指定的软件",
		Long: `启动指定的软件服务。

示例：
  servon software start nginx    # 启动 nginx
  servon software start mysql    # 启动 mysql`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				color.New(color.FgRed).Println("\n❌ 缺少软件名称参数")
				fmt.Println("\n用法:")
				color.New(color.FgYellow).Print("  servon software start ")
				fmt.Println("[软件名称]")

				// 显示支持的软件列表
				manager := softwares.NewSoftwareManager()
				names := manager.GetSoftwareNames()
				fmt.Println("\n支持的软件:")
				for _, name := range names {
					color.New(color.FgHiWhite).Printf("  - %s\n", name)
				}

				fmt.Println("\n示例:")
				color.New(color.FgCyan).Println("  servon software start nginx")
				color.New(color.FgCyan).Println("  servon software start mysql")
				return nil
			}

			manager := softwares.NewSoftwareManager()
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

			// 开始启动
			utils.InfoTitle("🚀 %s 启动中 ...", name)

			err := manager.StartSoftware(name, nil)
			if err != nil {
				utils.InfoTitle("❌ %s 启动失败", name)
				utils.Error("%s", err)
				return nil
			}

			utils.InfoTitle("✅ %s 启动成功！", name)

			return nil
		},
	}
}
