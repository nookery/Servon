package astro

import (
	"os"
)

func (a *AstroPlugin) build(path string) error {
	// 确保保存路径存在
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	// pnpm install
	if err := a.RunShell("pnpm", "install"); err != nil {
		return err
	}

	a.Info("pnpm install 成功")

	// pnpm build
	if err := a.RunShell("pnpm", "build"); err != nil {
		return err
	}

	a.Info("pnpm build 成功")

	return nil
}
