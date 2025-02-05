package deploy

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// newStartCmd 返回 start 子命令
func newStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "启动部署",
		Long:  `启动项目部署流程`,
		Run: func(cmd *cobra.Command, args []string) {
			color.New(color.FgHiCyan).Printf("Starting deployment...\n")
		},
	}
}
