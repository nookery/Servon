package serve

import (
	"servon/core"
)

type ServePlugin struct {
	*core.Core
}

func NewServePlugin(core *core.Core) *ServePlugin {
	return &ServePlugin{Core: core}
}

func Setup(core *core.Core) {
	plugin := NewServePlugin(core)

	core.AddCommand(plugin.NewServeCommand())
}
