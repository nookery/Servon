package astro

import (
	"path/filepath"
	"servon/core"
)

type AstroPlugin struct {
	*core.App
}

func Setup(app *core.App) {
	astro := NewAstroPlugin(app)

	// 添加 Astro 部署器到部署管理器
	app.AddDeployer(astro)
}

func NewAstroPlugin(app *core.App) *AstroPlugin {
	return &AstroPlugin{
		App: app,
	}
}

func (d *AstroPlugin) CanHandle(workDir string) bool {
	// 检查是否存在 astro.config.mjs 文件
	configPath := filepath.Join(workDir, "astro.config.mjs")
	return d.FileUtil.IsFileExists(configPath)
}

func (d *AstroPlugin) Build(workDir string, logger *core.LogUtil) error {
	return d.build(workDir)
}

func (d *AstroPlugin) Deploy(workDir string, targetDir string, logger *core.LogUtil) error {
	logger.Info("开始部署 Astro 项目")
	// 使用现有的 deploy 函数，但需要调整参数
	repo := filepath.Base(workDir) // 使用目录名作为项目名
	return d.deploy(repo, DefaultBranch, DefaultHost, DefaultPort)
}

func (d *AstroPlugin) GetName() string {
	return "astro"
}
