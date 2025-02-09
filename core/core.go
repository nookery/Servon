package core

import (
	"servon/core/api"
	"servon/core/model"
)

type Core struct {
	api.Command
	api.Shell
	api.Soft
	api.System
	api.Util
	api.Version
	api.Data
	api.Sample
}

type OSType = model.OSType

// New 创建Core实例
func New() *Core {
	return &Core{
		Command: api.NewCommand(),
		Shell:   api.NewShell(),
		Soft:    api.NewSoft(),
		System:  api.NewSystem(),
		Util:    api.NewUtil(),
		Version: api.NewVersion(),
		Data:    api.NewData(),
		Sample:  api.NewSample(),
	}
}
