package caddy

import (
	"servon/core"
)

func Setup(core *core.Core) {
	core.RegisterSoftware("caddy", NewCaddySoft(core))
	core.AddCommand(NewCaddyCommand(core))
}
