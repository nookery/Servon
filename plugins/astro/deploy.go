package astro

import (
	"fmt"

	"github.com/fatih/color"
)

func (a *AstroPlugin) deploy(repo string, storage string, port int) error {
	err := a.clone(repo, storage)
	if err != nil {
		return a.PrintAndReturnError(err.Error())
	}

	err = a.build(storage)
	if err != nil {
		return a.PrintAndReturnError(err.Error())
	}

	// 设置默认端口
	if port == 0 {
		port = 3000 // Astro 的默认端口
	}

	logChan := make(chan string)
	go func() {
		for msg := range logChan {
			fmt.Println(msg)
		}
	}()
	serviceFilePath, err := a.RunBackgroundService("node", []string{storage + "/index.js", "--port", fmt.Sprintf("%d", port)}, logChan)
	if err != nil {
		return a.PrintAndReturnError(err.Error())
	}

	// 成功提示
	fmt.Println()
	color.New(color.FgGreen, color.Bold).Printf("✨ Astro项目部署成功！\n")
	fmt.Println()
	color.New(color.FgWhite).Print("📦 仓库地址: ")
	color.New(color.FgHiWhite).Printf("%s\n", repo)
	color.New(color.FgWhite).Print("📁 存储路径: ")
	color.New(color.FgHiWhite).Printf("%s\n", storage)
	color.New(color.FgWhite).Print("📁 服务文件路径: ")
	color.New(color.FgHiWhite).Printf("%s\n", serviceFilePath)
	color.New(color.FgWhite).Print("🌐 服务端口: ")
	color.New(color.FgHiWhite).Printf("%d\n", port)
	fmt.Println()
	return nil
}
