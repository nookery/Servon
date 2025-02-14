package astro

import (
	"servon/core"
)

type AstroPlugin struct {
	*core.Core
}

func Setup(core *core.Core) {
	astro := NewAstroPlugin(core)

	core.AppendDeploySubCommand(astro.newAstroCommand())
}

func NewAstroPlugin(core *core.Core) *AstroPlugin {
	return &AstroPlugin{
		Core: core,
	}
}
