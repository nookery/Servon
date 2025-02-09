package core

import (
	"servon/core/provider"
)

type Core struct {
	softwareProvider provider.SoftwareProvider
	commandProvider  provider.CommandProvider
	versionProvider  provider.VersionProvider
}

// New 创建Core实例
func New() Core {
	return Core{
		softwareProvider: provider.NewSoftwareProvider(),
		commandProvider:  provider.NewCommandProvider(),
		versionProvider:  provider.NewVersionProvider(),
	}
}
