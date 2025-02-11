package caddy

import (
	"servon/core"
	"servon/core/contract"
)

func NewCaddySoft(core *core.Core) contract.SuperSoft {
	return &Caddy{
		Core: core,
		info: contract.SoftwareInfo{
			Name:        "caddy",
			Description: "现代化的 Web 服务器，支持自动 HTTPS",
		},
	}
}
