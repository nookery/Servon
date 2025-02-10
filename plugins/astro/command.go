package astro

import (
	"github.com/spf13/cobra"
	"servon/core"
)

func (a *AstroPlugin) newAstroCommand() *cobra.Command {
	rootCmd := a.NewCommand(core.CommandOptions{
		Use:   "astro",
		Short: "用来部署Astro项目",
		Run: func(cmd *cobra.Command, args []string) {
			repo, _ := cmd.Flags().GetString("repo")
			storage, _ := cmd.Flags().GetString("storage")
			port, _ := cmd.Flags().GetInt("port")

			err := a.deploy(repo, storage, port)
			if err != nil {
				a.Error(err.Error())
			}
		},
	})

	// 添加命令行参数
	rootCmd.Flags().String("repo", "", "Astro项目的Git仓库地址")
	rootCmd.Flags().String("storage", "", "项目部署的目标存储路径")
	rootCmd.Flags().Int("port", 0, "服务端口")
	rootCmd.MarkFlagRequired("repo")
	rootCmd.MarkFlagRequired("storage")

	return rootCmd
}
