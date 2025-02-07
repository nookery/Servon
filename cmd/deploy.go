package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"servon/cmd/deploy"
	"servon/cmd/internal/utils"
)

// DeployCmd 表示 deploy 命令
var DeployCmd = &cobra.Command{
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

func init() {
	DeployCmd.AddCommand(deployStartCmd)
	DeployCmd.AddCommand(deployStopCmd)
	DeployCmd.AddCommand(getProjectsCmd)
	DeployCmd.AddCommand(deployServeCmd)

	deployServeCmd.Flags().String("name", "", "服务名称")
	deployServeCmd.Flags().String("path", "", "本地文件夹路径")
	deployServeCmd.Flags().String("domain", "", "访问域名")
}

// deployStartCmd 表示 deploy start 子命令
var deployStartCmd = &cobra.Command{
	Use:   "start",
	Short: "启动部署",
	Long:  `启动项目部署流程`,
	Run: func(cmd *cobra.Command, args []string) {
		color.New(color.FgHiCyan).Printf("Starting deployment...\n")
	},
}

// deployStopCmd 表示 deploy stop 子命令
var deployStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止部署",
	Long:  `停止正在进行的部署流程`,
	Run: func(cmd *cobra.Command, args []string) {
		color.New(color.FgHiCyan).Printf("Stopping deployment...\n")
	},
}

var getProjectsCmd = &cobra.Command{
	Use:   "get",
	Short: "获取项目列表",
	Long:  `获取所有项目列表`,
	Run: func(cmd *cobra.Command, args []string) {
		projects, err := deploy.GetProjects()
		if err != nil {
			color.New(color.FgRed).Printf("获取项目列表失败: %v\n", err)
			return
		}
		color.New(color.FgGreen).Printf("项目列表: %v\n", projects)
	},
}

// deployServeCmd 表示 deploy serve 子命令
var deployServeCmd = &cobra.Command{
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

		err := deploy.ServeStatic(name, path, domain)
		if err != nil {
			color.New(color.FgRed).Printf("创建静态文件服务失败: %v\n", err)
			return
		}

		// 成功提示
		fmt.Println() // 添加空行使显示更清晰
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
