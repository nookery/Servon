package api

import (
	"servon/core/libs"
)

type Version struct {
	versionProvider libs.VersionProvider
}

func NewVersion() Version {
	return Version{
		versionProvider: libs.NewVersionProvider(),
	}
}

func (c *Version) GetVersion() string {
	return c.versionProvider.GetVersion()
}

func (c *Version) GetFullVersionInfo() string {
	return c.versionProvider.GetFullVersionInfo()
}
