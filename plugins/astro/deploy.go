package astro

import (
	"fmt"
	"servon/core"

	"github.com/fatih/color"
)

func deploy(core *core.Core, repo string, storage string) error {
	err := clone(core, repo, storage)
	if err != nil {
		return core.PrintAndReturnError(err.Error())
	}

	err = build(core, storage)
	if err != nil {
		return core.PrintAndReturnError(err.Error())
	}

	// 成功提示
	fmt.Println()
	color.New(color.FgGreen, color.Bold).Printf("✨ Astro项目部署成功！\n")
	fmt.Println()
	color.New(color.FgWhite).Print("📦 仓库地址: ")
	color.New(color.FgHiWhite).Printf("%s\n", repo)
	color.New(color.FgWhite).Print("📁 存储路径: ")
	color.New(color.FgHiWhite).Printf("%s\n", storage)
	fmt.Println()
	return nil
}
