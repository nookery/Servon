package api

import (
	"fmt"
	"servon/core/contract"
	"servon/core/libs"
)

type Soft struct {
	Softwares map[string]contract.SuperSoft
}

func NewSoft() Soft {
	return Soft{
		Softwares: make(map[string]contract.SuperSoft),
	}
}

// Install 安装软件, 如果提供了日志通道则输出日志
func (c *Soft) Install(name string, logChan chan<- string) error {
	software, ok := c.Softwares[name]
	if !ok {
		registeredSoftwares := make([]string, 0, len(c.Softwares))
		for name := range c.Softwares {
			registeredSoftwares = append(registeredSoftwares, name)
		}
		return libs.PrintAndReturnError(fmt.Sprintf("软件 %s 未注册, 可用的软件有: %v", name, registeredSoftwares))
	}
	return software.Install(logChan)
}

// UninstallSoftware 卸载软件
func (c *Soft) UninstallSoftware(name string, logChan chan<- string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return libs.PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}
	return software.Uninstall(logChan)
}

// StartSoftware 启动软件
func (c *Soft) StartSoftware(name string, logChan chan<- string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return libs.PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}
	return software.Start(logChan)
}

// StopSoftware 停止软件
func (c *Soft) StopSoftware(name string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return libs.PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}
	return software.Stop()
}

// GetSoftwareStatus 获取软件状态
func (c *Soft) GetSoftwareStatus(name string) (map[string]string, error) {
	software, ok := c.Softwares[name]
	if !ok {
		return nil, libs.PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}
	return software.GetStatus()
}

// RegisterSoftware 注册软件
func (c *Soft) RegisterSoftware(name string, software contract.SuperSoft) error {
	if _, exists := c.Softwares[name]; exists {
		return libs.PrintAndReturnError(fmt.Sprintf("软件 %s 已注册", name))
	}
	c.Softwares[name] = software
	return nil
}

// GetAllSoftware 获取所有软件
func (c *Soft) GetAllSoftware() []string {
	softwareNames := make([]string, 0, len(c.Softwares))
	for name := range c.Softwares {
		softwareNames = append(softwareNames, name)
	}
	return softwareNames
}
