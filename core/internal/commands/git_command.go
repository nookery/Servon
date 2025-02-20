package commands

import (
	"servon/core/internal/managers"

	"github.com/spf13/cobra"
)

// GetGitRootCommand 获取git命令根命令，返回一个 cobra.Command
func GetGitRootCommand(g *managers.GitManager) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "git",
		Short: "git 命令，比默认的 git 命令更智能",
	}

	rootCmd.AddCommand(GetCloneCommand(g))

	return rootCmd
}

// GetCloneCommand 获取克隆命令，返回一个 cobra.Command
func GetCloneCommand(g *managers.GitManager) *cobra.Command {
	cloneCmd := &cobra.Command{
		Use:   "clone",
		Short: "克隆一个git仓库",
		Run: func(cmd *cobra.Command, args []string) {
			url, _ := cmd.Flags().GetString("url")
			branch, _ := cmd.Flags().GetString("branch")
			targetDir, _ := cmd.Flags().GetString("target-dir")

			err := g.GitClone(url, branch, targetDir)
			if err != nil {
				PrintError(err)
			}
		},
	}

	cloneCmd.Flags().StringP("url", "u", "", "仓库地址")
	cloneCmd.Flags().StringP("branch", "b", "master", "分支")
	cloneCmd.Flags().StringP("target-dir", "t", "", "目标目录")

	return cloneCmd
}
