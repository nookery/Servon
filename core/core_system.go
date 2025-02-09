package core

import "servon/core/model"

func (c *Core) InstallSoftwareUseSystem(name string) error {
	return c.systemProvider.InstallSoftware(name, nil)
}

func (c *Core) UninstallSoftwareUseSystem(name string) error {
	return c.systemProvider.UninstallSoftware(name, nil)
}

func (c *Core) GetOSType() model.OSType {
	return c.systemProvider.GetOSType()
}

func (c *Core) CanUseApt() bool {
	return c.systemProvider.CanUseApt()
}
