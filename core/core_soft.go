package core

import "servon/core/provider"

// InstallSoftware 安装软件
func (c *Core) InstallSoftware(name string, logChan chan<- string) error {
	return c.softwareProvider.Install(name, logChan)
}

// UninstallSoftware 卸载软件
func (c *Core) UninstallSoftware(name string, logChan chan<- string) error {
	return c.softwareProvider.Uninstall(name, logChan)
}

// StartSoftware 启动软件
func (c *Core) StartSoftware(name string, logChan chan<- string) error {
	return c.softwareProvider.Start(name, logChan)
}

// StopSoftware 停止软件
func (c *Core) StopSoftware(name string) error {
	return c.softwareProvider.Stop(name)
}

// GetSoftwareStatus 获取软件状态
func (c *Core) GetSoftwareStatus(name string) (map[string]string, error) {
	return c.softwareProvider.GetStatus(name)
}

// RegisterSoftware 注册软件
func (c *Core) RegisterSoftware(name string, software provider.SuperSoft) error {
	return c.softwareProvider.Register(name, software)
}

// GetAllSoftware 获取所有软件
func (c *Core) GetAllSoftware() []string {
	return c.softwareProvider.GetAllSoftware()
}
