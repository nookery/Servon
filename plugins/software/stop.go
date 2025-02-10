package software

import (
	"fmt"
	"servon/core"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func (p *SoftWarePlugin) newStopCmd() *cobra.Command {
	return p.core.NewCommand(core.CommandOptions{
		Use:   "stop [软件名称]",
		Short: "停止指定的软件",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				color.New(color.FgRed).Println("\n❌ 缺少软件名称参数")
				fmt.Println("\n用法:")
				color.New(color.FgYellow).Print("  servon software stop ")
				fmt.Println("[软件名称]")

				// 显示支持的软件列表
				names := p.core.GetAllSoftware()
				fmt.Println("\n支持的软件:")
				for _, name := range names {
					color.New(color.FgHiWhite).Printf("  - %s\n", name)
				}

				fmt.Println("\n示例:")
				color.New(color.FgCyan).Println("  servon software stop caddy")
				color.New(color.FgCyan).Println("  servon software stop clash")
				return
			}

			name := args[0]

			// 检查软件是否支持
			supported := false
			for _, sw := range p.core.GetAllSoftware() {
				if sw == name {
					supported = true
					break
				}
			}

			if !supported {
				color.New(color.FgRed).Printf("\n❌ 不支持的软件: %s\n", name)
				fmt.Println("\n支持的软件:")
				for _, sw := range p.core.GetAllSoftware() {
					color.New(color.FgHiWhite).Printf("  - %s\n", sw)
				}
				return
			}

			// 开始停止
			p.core.Infoln("🛑 %s 停止中 ...", name)

			err := p.core.StopSoftware(name)
			if err != nil {
				p.core.Infoln("❌ %s 停止失败", name)
				p.core.Error("%s", err)
				return
			}

			p.core.Infoln("✅ %s 已停止！", name)
		},
	})
}
