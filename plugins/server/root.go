package server

import (
	"servon/components/logger"
	"servon/components/string_util"
	"servon/components/web_server"
	"servon/core"
	"servon/plugins/server/web/routers"

	"github.com/spf13/cobra"
)

var log = logger.DefaultLogUtil
var server = web_server.NewWebServer()
var stringUtil = string_util.DefaultStringUtil

// Setup 注册port插件到应用程序
func Setup(app *core.App) {
	manager := app.FullManager
	routers.Setup(manager, server.Engine, true)
	app.GetRootCommand().AddCommand(GetServerCommand())
}

// GetServerCommand 获取服务器管理命令
func GetServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "启动服务器",
	}

	cmd.AddCommand(MakeStartCommand())
	cmd.AddCommand(MakeStopCommand())
	cmd.AddCommand(MakeRestartCommand())

	return cmd
}
