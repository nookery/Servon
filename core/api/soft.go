package api

import "servon/core/libs"

type Soft struct {
	softwareProvider libs.SoftwareProvider
}

func NewSoft() Soft {
	return Soft{
		softwareProvider: libs.NewSoftwareProvider(),
	}
}

// InstallSoftware 安装软件
func (c *Soft) InstallSoftware(name string, logChan chan<- string) error {
	return c.softwareProvider.Install(name, logChan)
}

// UninstallSoftware 卸载软件
func (c *Soft) UninstallSoftware(name string, logChan chan<- string) error {
	return c.softwareProvider.Uninstall(name, logChan)
}

// StartSoftware 启动软件
func (c *Soft) StartSoftware(name string, logChan chan<- string) error {
	return c.softwareProvider.Start(name, logChan)
}

// StopSoftware 停止软件
func (c *Soft) StopSoftware(name string) error {
	return c.softwareProvider.Stop(name)
}

// GetSoftwareStatus 获取软件状态
func (c *Soft) GetSoftwareStatus(name string) (map[string]string, error) {
	return c.softwareProvider.GetStatus(name)
}

// RegisterSoftware 注册软件
func (c *Soft) RegisterSoftware(name string, software libs.SuperSoft) error {
	return c.softwareProvider.Register(name, software)
}

// GetAllSoftware 获取所有软件
func (c *Soft) GetAllSoftware() []string {
	return c.softwareProvider.GetAllSoftware()
}
