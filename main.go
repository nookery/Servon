package main

import (
	"os"
	"servon/core"

	"github.com/fatih/color"

	// 导入插件
	"servon/plugins/caddy"
	"servon/plugins/clash"
	"servon/plugins/deploy"
	"servon/plugins/git"
	"servon/plugins/ip"
	"servon/plugins/nodejs"
	"servon/plugins/pnpm"
	"servon/plugins/serve"
	"servon/plugins/software"
	"servon/plugins/version"
	"servon/plugins/yarn"
)

func main() {
	core := core.New()

	serve.Setup(&core)
	software.Setup(&core)
	caddy.Setup(&core)
	nodejs.Setup(&core)
	yarn.Setup(&core)
	git.Setup(&core)
	pnpm.Setup(&core)
	clash.Setup(&core)
	ip.Setup(&core)
	version.Setup(&core)
	deploy.Setup(&core)
	if err := core.GetRootCommand().Execute(); err != nil {
		color.Red("Error: %v\n", err)
		os.Exit(1)
	}
}
