package software

import (
	"fmt"
	"servon/cmd/internal/softwares"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// newInfoCmd 返回 info 子命令
func newInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info [软件名称]",
		Short: "显示软件详细信息",
		Long: `显示指定软件的详细信息，包括安装状态、版本等。

示例：
  servon software info nginx    # 显示 nginx 的详细信息
  servon software info mysql    # 显示 mysql 的详细信息`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				color.New(color.FgRed).Println("\n❌ 缺少软件名称参数")
				fmt.Println("\n用法:")
				color.New(color.FgYellow).Print("  servon software info ")
				fmt.Println("[软件名称]")

				// 显示支持的软件列表
				manager := softwares.NewSoftwareManager()
				names := manager.GetSoftwareNames()
				fmt.Println("\n支持的软件:")
				for _, name := range names {
					color.New(color.FgHiWhite).Printf("  - %s\n", name)
				}

				fmt.Println("\n示例:")
				color.New(color.FgCyan).Println("  servon software info nginx")
				color.New(color.FgCyan).Println("  servon software info mysql")
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

			// 获取软件状态
			status, err := manager.GetSoftwareStatus(name)
			if err != nil {
				color.New(color.FgRed).Printf("\n❌ 获取软件状态失败: %v\n", err)
				return nil
			}

			// 显示软件信息
			fmt.Println()
			color.New(color.FgCyan, color.Bold).Printf("📦 %s\n", name)
			fmt.Println()

			// 显示安装状态
			color.New(color.FgWhite).Print("状态: ")
			switch status["status"] {
			case "running":
				color.New(color.FgGreen).Println("运行中")
			case "stopped":
				color.New(color.FgYellow).Println("已停止")
			case "not_installed":
				color.New(color.FgRed).Println("未安装")
			default:
				color.New(color.FgHiWhite).Println(status["status"])
			}

			// 显示版本信息
			if version := status["version"]; version != "" {
				color.New(color.FgWhite).Print("版本: ")
				color.New(color.FgHiWhite).Println(version)
			}

			fmt.Println()
			return nil
		},
	}
}
