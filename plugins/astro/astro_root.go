package astro

import (
	"servon/core"
)

type AstroPlugin struct {
	*core.App
}

func Setup(app *core.App) {
	astro := NewAstroPlugin(app)

	app.AppendDeploySubCommand(astro.newAstroCommand())
}

func NewAstroPlugin(app *core.App) *AstroPlugin {
	return &AstroPlugin{
		App: app,
	}
}
