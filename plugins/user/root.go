package user

import (
	"servon/core"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// UserCmd represents the user command
var UserCmd = &cobra.Command{
	Use:   "user",
	Short: color.Blue.Render("用户管理工具"),
	Long:  color.Success.Render("\r\n用户管理工具，用于列出系统用户、创建用户、删除用户等操作"),
}

func init() {
	// 添加子命令
	UserCmd.AddCommand(listCmd)
	UserCmd.AddCommand(createCmd)
	UserCmd.AddCommand(deleteCmd)
	UserCmd.AddCommand(infoCmd)
}

// Setup 注册用户管理插件到应用程序
func Setup(app *core.App) {
	app.GetRootCommand().AddCommand(UserCmd)
}
