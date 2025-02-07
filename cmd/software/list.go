package software

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// newListCmd 返回 list 子命令
func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "显示支持的软件列表",
		Run: func(cmd *cobra.Command, args []string) {
			manager := NewSoftwareManager()
			names := manager.GetSoftwareNames()

			color.New(color.FgHiCyan).Println("\n支持的软件列表：")
			for _, name := range names {
				color.New(color.FgHiWhite).Printf("- %s\n", name)
			}
			color.New(color.FgHiCyan).Println()
		},
	}
}
