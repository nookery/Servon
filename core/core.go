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
	api.VersionApi
	api.Data
	api.LogApi
	api.DeployApi
	*libs.Printer
	*libs.PortManager
	*libs.BasicInfoManager
	*libs.OSInfoManager
	*libs.SystemResourcesManager
	*libs.ProcessManager
	*libs.FilesManager
	*libs.NetworkManager
	*libs.ServiceManager
	*libs.AptManager
	*libs.Dpkg
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
		CommandApi:             api.NewCommandApi(),
		Soft:                   api.NewSoft(),
		VersionApi:             api.NewVersion(),
		Data:                   api.NewData(),
		LogApi:                 api.NewLogApi(),
		DeployApi:              api.NewDeployApi(),
		Printer:                libs.NewPrinter(),
		PortManager:            libs.NewPortManager(),
		BasicInfoManager:       libs.NewBasicInfoManager(),
		OSInfoManager:          libs.NewOSInfoManager(),
		SystemResourcesManager: libs.NewSystemResourcesManager(),
		ProcessManager:         libs.NewProcessManager(),
		FilesManager:           libs.NewFilesManager(),
		NetworkManager:         libs.NewNetworkManager(),
		ServiceManager:         libs.NewServiceManager(),
		AptManager:             libs.NewAptManager(),
		Dpkg:                   libs.NewDpkg(),
	}

	core.AddCommand(core.GetDeployCommand())
	core.AddCommand(core.GetVersionCommand())
	core.AddCommand(core.GetSoftwareCommand())

	return core
}
