package serve

import (
	"servon/core"
)

func Setup(core *core.Core) {
	core.AddCommand(ServeCmd)
}
