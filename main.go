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
	"servon/plugins/nodejs"
	"servon/plugins/npm"
	"servon/plugins/pnpm"
	"servon/plugins/yarn"
)

func main() {
	core := core.New()

	if core.CompareEnv() != nil {
		color.Red("Error: .env 文件与 .env.example 文件不一致")
		os.Exit(1)
	}

	core.LoadEnv()

	caddy.Setup(core)
	nodejs.Setup(core)
	yarn.Setup(core)
	git.Setup(core)
	pnpm.Setup(core)
	clash.Setup(core)
	ip.Setup(core)
	astro.Setup(core)
	npm.Setup(core)
	github_runner.Setup(core)

	if err := core.GetRootCommand().Execute(); err != nil {
		color.Red("Error: %v\n", err)
		os.Exit(1)
	}
}
