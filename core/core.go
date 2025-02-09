package core

import (
	"servon/core/model"
	"servon/core/provider"
)

type Core struct {
	softwareProvider provider.SoftwareProvider
	commandProvider  provider.CommandProvider
	versionProvider  provider.VersionProvider
	systemProvider   provider.SystemProvider
	sampleProvider   provider.SampleProvider
	shellProvider    provider.ShellProvider
	utilProvider     provider.UtilProvider
}

type OSType = model.OSType

// New 创建Core实例
func New() *Core {
	return &Core{
		softwareProvider: provider.NewSoftwareProvider(),
		commandProvider:  provider.NewCommandProvider(),
		versionProvider:  provider.NewVersionProvider(),
		systemProvider:   provider.NewSystemProvider(),
		sampleProvider:   provider.NewSampleProvider(),
		shellProvider:    provider.NewShellProvider(),
		utilProvider:     provider.NewUtilProvider(),
	}
}
