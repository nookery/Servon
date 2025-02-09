package deploy

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// newGetCmd 返回 get 子命令
func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "获取项目列表",
		Long:  `获取所有项目列表`,
		Run: func(cmd *cobra.Command, args []string) {
			projects, err := GetProjects()
			if err != nil {
				color.New(color.FgRed).Printf("获取项目列表失败: %v\n", err)
				return
			}
			color.New(color.FgGreen).Printf("项目列表: %v\n", projects)
		},
	}
}
