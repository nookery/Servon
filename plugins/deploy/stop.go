package deploy

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// newStopCmd 返回 stop 子命令
func newStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "停止部署",
		Long:  `停止正在进行的部署流程`,
		Run: func(cmd *cobra.Command, args []string) {
			color.New(color.FgHiCyan).Printf("Stopping deployment...\n")
		},
	}
}
