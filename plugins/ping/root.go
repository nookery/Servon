package ping

import (
	"fmt"
	"servon/core"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// PingCmd represents the ping command
var PingCmd = &cobra.Command{
	Use:   "ping",
	Short: color.Blue.Render("输出pang"),
	Long:  color.Success.Render(`输出pang，用于测试。`),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pang")
	},
}

// Setup 注册ping插件到应用程序
func Setup(app *core.App) {
	app.GetRootCommand().AddCommand(PingCmd)
}
