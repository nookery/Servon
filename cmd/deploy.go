package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"servon/internal/deploy"
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
		cmd.Println("使用方法:\tservon deploy COMMAND")
		cmd.Println("\nServon 部署管理命令\n")
		cmd.Println("可用命令:")
		cmd.Println("  start\t启动部署")
		cmd.Println("  stop\t停止部署")
		cmd.Println("\n使用 \"servon deploy COMMAND --help\" 了解更多命令信息")
		return nil
	},
}

func init() {
	DeployCmd.AddCommand(deployStartCmd)
	DeployCmd.AddCommand(deployStopCmd)
	DeployCmd.AddCommand(getProjectsCmd)
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
