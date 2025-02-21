package astro

import (
	"servon/core"

	"github.com/spf13/cobra"
)

func (a *AstroPlugin) newAstroCommand() *cobra.Command {
	rootCmd := a.NewCommand(core.CommandOptions{
		Use:   "astro",
		Short: "用来部署Astro项目",
		Run: func(cmd *cobra.Command, args []string) {
			repo, _ := cmd.Flags().GetString("repo")
			port, _ := cmd.Flags().GetInt("port")
			branch, _ := cmd.Flags().GetString("branch")
			host, _ := cmd.Flags().GetString("host")

			err := a.deploy(repo, branch, host, port)
			if err != nil {
				a.Error(err)
			}
		},
	})

	// 添加命令行参数
	rootCmd.Flags().String("repo", "", "Astro项目的Git仓库地址")
	rootCmd.Flags().String("branch", "main", "Astro项目的分支")
	rootCmd.Flags().Int("port", 0, "服务端口")
	rootCmd.Flags().String("host", "0.0.0.0", "服务Host")
	rootCmd.MarkFlagRequired("repo")

	return rootCmd
}
