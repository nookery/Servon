package astro

import (
	"fmt"
	"servon/core"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func Setup(core *core.Core) {
	core.AddCommand(NewAstroCommand(core))
}

func NewAstroCommand(core *core.Core) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "astro",
		Short: "用来部署Astro项目",
		Long: `部署Astro项目到指定目录。

示例：
  servon astro --repo https://github.com/user/project --storage /path/to/storage`,
		Run: func(cmd *cobra.Command, args []string) {
			repo, _ := cmd.Flags().GetString("repo")
			storage, _ := cmd.Flags().GetString("storage")

			if repo == "" || storage == "" {
				color.New(color.FgRed).Println("\n❌ 缺少必要参数")
				fmt.Println("\n必需参数:")
				if repo == "" {
					color.New(color.FgYellow).Print("  --repo    ")
					fmt.Println("Astro项目的Git仓库地址")
				}
				if storage == "" {
					color.New(color.FgYellow).Print("  --storage ")
					fmt.Println("项目部署的目标存储路径")
				}

				fmt.Println("\n示例:")
				color.New(color.FgCyan).Println("  servon astro \\")
				color.New(color.FgCyan).Println("    --repo https://github.com/user/project \\")
				color.New(color.FgCyan).Println("    --storage /path/to/storage")
				return
			}

			err := deploy(core, repo, storage)
			if err != nil {
				core.PrintCommandErrorAndExit(err)
			}
		},
	}

	// 添加命令行参数
	rootCmd.Flags().String("repo", "", "Astro项目的Git仓库地址")
	rootCmd.Flags().String("storage", "", "项目部署的目标存储路径")

	return rootCmd
}
