package core

func (c *Core) GetVersion() string {
	return c.versionProvider.GetVersion()
}

func (c *Core) GetFullVersionInfo() string {
	return c.versionProvider.GetFullVersionInfo()
}
