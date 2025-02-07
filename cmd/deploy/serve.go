package deploy

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// newServeCmd 返回 serve 子命令
func newServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "创建静态文件服务",
		Long: `将本地文件夹暴露为Web服务。

示例：
  servon deploy serve --name myfiles --path /path/to/files --domain files.example.com`,
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			path, _ := cmd.Flags().GetString("path")
			domain, _ := cmd.Flags().GetString("domain")

			if name == "" || path == "" || domain == "" {
				color.New(color.FgRed).Println("\n❌ 缺少必要参数")
				fmt.Println("\n必需参数:")
				if name == "" {
					color.New(color.FgYellow).Print("  --name  ")
					fmt.Println("服务名称，用于标识和管理此服务")
				}
				if path == "" {
					color.New(color.FgYellow).Print("  --path  ")
					fmt.Println("本地文件夹路径，指定要分享的文件目录")
				}
				if domain == "" {
					color.New(color.FgYellow).Print("  --domain")
					fmt.Println("访问域名，用于通过浏览器访问文件")
				}

				fmt.Println("\n示例:")
				color.New(color.FgCyan).Println("  servon deploy serve \\")
				color.New(color.FgCyan).Println("    --name photos \\")
				color.New(color.FgCyan).Println("    --path /home/user/pictures \\")
				color.New(color.FgCyan).Println("    --domain photos.example.com")
				return
			}

			err := ServeStatic(name, path, domain)
			if err != nil {
				color.New(color.FgRed).Printf("创建静态文件服务失败: %v\n", err)
				return
			}

			// 成功提示
			fmt.Println()
			color.New(color.FgGreen, color.Bold).Printf("✨ 静态文件服务创建成功！\n")
			fmt.Println()
			color.New(color.FgWhite).Print("📂 服务名称: ")
			color.New(color.FgHiWhite).Printf("%s\n", name)
			color.New(color.FgWhite).Print("📁 文件路径: ")
			color.New(color.FgHiWhite).Printf("%s\n", path)
			color.New(color.FgWhite).Print("🌐 访问地址: ")
			color.New(color.FgHiCyan).Printf("http://%s\n", domain)
			fmt.Println()
			color.New(color.FgHiBlack).Println("提示：请确保域名已正确解析到服务器IP")
		},
	}
}
