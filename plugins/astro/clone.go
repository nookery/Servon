package astro

import (
	"fmt"
	"os"
	"servon/core"
)

func clone(core *core.Core, address string, savePath string) error {
	// 如果路径存在，且不为空，则不能克隆
	if _, err := os.Stat(savePath); err == nil {
		return core.PrintCommandErrorAndExit(fmt.Errorf("路径 %s 已存在，不能克隆", savePath))
	}

	// 检查并输出代理配置
	if out, err := core.RunShellWithOutput("git", "config", "--global", "--get", "http.proxy"); err == nil && out != "" {
		fmt.Printf("Git HTTP 代理: %s", out)
	}
	if out, err := core.RunShellWithOutput("git", "config", "--global", "--get", "https.proxy"); err == nil && out != "" {
		fmt.Printf("Git HTTPS 代理: %s", out)
	}

	// 执行git clone命令
	if err := core.RunShell("git", "clone", address, savePath); err != nil {
		return err
	}

	core.PrintStepFinish("Astro项目克隆成功")
	return nil
}
