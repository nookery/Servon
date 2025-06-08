package main

import (
	"os"
	"servon/core"

	"github.com/fatih/color"

	// 导入插件
	"servon/plugins/astro"
	"servon/plugins/caddy"
	"servon/plugins/clash"
	"servon/plugins/git"
	"servon/plugins/github_runner"
	"servon/plugins/ip"
	"servon/plugins/joke"
	"servon/plugins/nodejs"
	"servon/plugins/npm"
	"servon/plugins/ping"
	"servon/plugins/pm2"
	"servon/plugins/pnpm"
	"servon/plugins/port"
	"servon/plugins/supervisor"
	"servon/plugins/xcode"
	"servon/plugins/yarn"
)

func main() {
	app := core.New()

	// 注册插件
	astro.Setup(app)
	caddy.Setup(app)
	clash.Setup(app)
	git.Setup(app)
	github_runner.Setup(app)
	ip.Setup(app)
	joke.Setup(app)
	nodejs.Setup(app)
	npm.Setup(app)
	ping.Setup(app)
	pm2.Setup(app)
	pnpm.Setup(app)
	port.Setup(app)
	supervisor.Setup(app)
	xcode.Setup(app)
	yarn.Setup(app)

	if err := app.GetRootCommand().Execute(); err != nil {
		color.Red("Error: %v\n", err)
		os.Exit(1)
	}
}
