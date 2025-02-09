package version

import (
	"servon/core"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func Setup(core *core.Core) {
	core.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Long:  "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			color.Green("Servon 版本: %s", core.GetVersion())
		},
	})
}
