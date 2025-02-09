package main

import (
	"os"
	"servon/core"

	"github.com/fatih/color"

	// 导入插件
	"servon/plugins/serve"
	"servon/plugins/software"
)

func main() {
	core := core.New()

	serve.Setup(&core)
	software.Setup(&core)

	if err := core.GetRootCommand().Execute(); err != nil {
		color.Red("Error: %v\n", err)
		os.Exit(1)
	}
}
