package software

import (
	"fmt"
	"servon/core"
	"servon/core/utils/logger"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func newStopCmd(core *core.Core) *cobra.Command {
	return &cobra.Command{
		Use:   "stop [软件名称]",
		Short: "停止指定的软件",
		Long: `停止指定的软件。

示例：
  servon software stop caddy    # 停止 Caddy 服务
  servon software stop clash    # 停止 Clash 服务`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				color.New(color.FgRed).Println("\n❌ 缺少软件名称参数")
				fmt.Println("\n用法:")
				color.New(color.FgYellow).Print("  servon software stop ")
				fmt.Println("[软件名称]")

				// 显示支持的软件列表
				names := core.GetAllSoftware()
				fmt.Println("\n支持的软件:")
				for _, name := range names {
					color.New(color.FgHiWhite).Printf("  - %s\n", name)
				}

				fmt.Println("\n示例:")
				color.New(color.FgCyan).Println("  servon software stop caddy")
				color.New(color.FgCyan).Println("  servon software stop clash")
				return nil
			}

			name := args[0]

			// 检查软件是否支持
			supported := false
			for _, sw := range core.GetAllSoftware() {
				if sw == name {
					supported = true
					break
				}
			}

			if !supported {
				color.New(color.FgRed).Printf("\n❌ 不支持的软件: %s\n", name)
				fmt.Println("\n支持的软件:")
				for _, sw := range core.GetAllSoftware() {
					color.New(color.FgHiWhite).Printf("  - %s\n", sw)
				}
				return nil
			}

			// 开始停止
			logger.InfoTitle("🛑 %s 停止中 ...", name)

			err := core.StopSoftware(name)
			if err != nil {
				logger.InfoTitle("❌ %s 停止失败", name)
				logger.Error("%s", err)
				return nil
			}

			logger.InfoTitle("✅ %s 已停止！", name)

			return nil
		},
	}
}
