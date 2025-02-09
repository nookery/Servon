package astro

import (
	"os"
	"servon/core"
)

func build(core *core.Core, path string) error {
	// 确保保存路径存在
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	// pnpm install
	if err := core.RunShell("pnpm", "install"); err != nil {
		return err
	}

	core.PrintStepFinish("pnpm install 成功")

	// pnpm build
	if err := core.RunShell("pnpm", "build"); err != nil {
		return err
	}

	core.PrintStepFinish("pnpm build 成功")

	return nil
}
