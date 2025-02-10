package core

import (
	"servon/core/api"
	"servon/core/libs"
)

// 调用关系 Core -> Core API(Public) -> Libs(Private)

const DataRootFolder = "/data"
const LoggerFolder = "/logs"

type Core struct {
	api.CommandApi
	api.Soft
	api.SystemApi
	api.VersionApi
	api.Data
	api.LogApi
	api.DeployApi
	api.PrintApi
}

type OSType = libs.OSType
type CommandOptions = api.CommandOptions

const (
	Ubuntu  OSType = "ubuntu"
	Debian  OSType = "debian"
	CentOS  OSType = "centos"
	RedHat  OSType = "redhat"
	Unknown OSType = "unknown"
)

// New 创建Core实例
func New() *Core {
	core := &Core{
		CommandApi: api.NewCommandApi(),
		Soft:       api.NewSoft(),
		SystemApi:  api.NewSystemApi(),
		VersionApi: api.NewVersion(),
		Data:       api.NewData(),
		LogApi:     api.NewLogApi(),
		DeployApi:  api.NewDeployApi(),
		PrintApi:   api.NewPrintApi(),
	}

	core.AddCommand(core.GetDeployCommand())
	core.AddCommand(core.GetVersionCommand())
	core.AddCommand(core.GetSoftwareCommand())

	return core
}
