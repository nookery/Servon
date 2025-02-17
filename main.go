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
	"servon/plugins/nodejs"
	"servon/plugins/npm"
	"servon/plugins/pm2"
	"servon/plugins/pnpm"
	"servon/plugins/supervisor"
	"servon/plugins/web"
	"servon/plugins/yarn"
)

func main() {
	app := core.New()

	caddy.Setup(app)
	nodejs.Setup(app)
	yarn.Setup(app)
	git.Setup(app)
	pnpm.Setup(app)
	clash.Setup(app)
	astro.Setup(app)
	npm.Setup(app)
	github_runner.Setup(app)
	web.Setup(app)
	pm2.Setup(app)
	supervisor.Setup(app)

	if err := app.GetRootCommand().Execute(); err != nil {
		color.Red("Error: %v\n", err)
		os.Exit(1)
	}
}
