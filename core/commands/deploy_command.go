package commands

import (
	"fmt"
	"servon/core/managers"

	"github.com/spf13/cobra"
)

// deployRootCommand 部署根命令
func makeDeployCommand(manager *managers.FullManager) *cobra.Command {
	return NewCommand(CommandOptions{
		Use:     "deploy [repository]",
		Short:   "部署项目",
		Aliases: []string{"d"},
		Run: func(cmd *cobra.Command, args []string) {
			deployManager := manager.DeployManager

			// 如果没有提供参数，显示帮助信息
			if len(args) == 0 {
				fmt.Println("使用方法: servon deploy [repository]")
				fmt.Println("\n仓库地址格式:")
				fmt.Println("  - 完整URL:   https://github.com/user/repo")
				fmt.Println("  - 简短格式:   user/repo")
				fmt.Println("\n支持的项目类型:")

				// 获取所有已注册的部署器
				deployers := deployManager.GetDeployers()
				if len(deployers) == 0 {
					fmt.Println("  当前未注册任何部署器")
				} else {
					for _, d := range deployers {
						fmt.Printf("  - %s\n", d.GetName())
					}
				}

				fmt.Println("\n示例:")
				fmt.Println("  servon deploy username/project")
				fmt.Println("  servon deploy https://github.com/username/project")
				return
			}

			deployManager.DeployProject(args[0])
		},
	})
}

func GetDeployCommand(manager *managers.FullManager) *cobra.Command {
	return makeDeployCommand(manager)
}
