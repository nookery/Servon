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
	"servon/plugins/ip"
	"servon/plugins/nodejs"
	"servon/plugins/pnpm"
	"servon/plugins/version"
	"servon/plugins/yarn"
)

func main() {
	core := core.New()

	caddy.Setup(core)
	nodejs.Setup(core)
	yarn.Setup(core)
	git.Setup(core)
	pnpm.Setup(core)
	clash.Setup(core)
	ip.Setup(core)
	version.Setup(core)
	astro.Setup(core)
	if err := core.GetRootCommand().Execute(); err != nil {
		color.Red("Error: %v\n", err)
		os.Exit(1)
	}
}
