package api

import (
	"servon/core/provider"
)

type Version struct {
	versionProvider provider.VersionProvider
}

func NewVersion() Version {
	return Version{
		versionProvider: provider.NewVersionProvider(),
	}
}

func (c *Version) GetVersion() string {
	return c.versionProvider.GetVersion()
}

func (c *Version) GetFullVersionInfo() string {
	return c.versionProvider.GetFullVersionInfo()
}
