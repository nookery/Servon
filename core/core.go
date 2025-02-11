package core

import (
	"servon/core/api"
	"servon/core/libs"
	"servon/core/serve"
)

// 调用关系 Core -> Core API -> Libs
// 或 Core -> Libs

const DataRootFolder = "/data"
const LoggerFolder = "/logs"

type Core struct {
	api.CommandApi
	api.Data
	api.LogApi
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
	*libs.CronManager
	*libs.VersionManager
	*libs.SoftManager
	*libs.DeployManager
}

type OSType = libs.OSType
type CommandOptions = api.CommandOptions
type CronTask = libs.CronTask
type ValidationError = libs.ValidationError
type ValidationErrors = libs.ValidationErrors

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
		SoftManager:            libs.DefaultSoftManager,
		Data:                   api.NewData(),
		LogApi:                 api.NewLogApi(),
		DeployManager:          libs.NewDeployManager(),
		Printer:                libs.NewPrinter(),
		PortManager:            libs.NewPortManager(),
		BasicInfoManager:       libs.NewBasicInfoManager(),
		OSInfoManager:          libs.NewOSInfoManager(),
		SystemResourcesManager: libs.DefaultSystemResourcesManager,
		ProcessManager:         libs.NewProcessManager(),
		FilesManager:           libs.NewFilesManager(),
		NetworkManager:         libs.NewNetworkManager(),
		ServiceManager:         libs.NewServiceManager(),
		AptManager:             libs.NewAptManager(),
		Dpkg:                   libs.NewDpkg(),
		CronManager:            libs.NewCronManager(),
		VersionManager:         libs.DefaultVersionManager,
	}

	core.AddCommand(core.GetDeployCommand())
	core.AddCommand(core.GetVersionCommand())
	core.AddCommand(core.GetSoftwareCommand())
	core.AddCommand(serve.NewServeCommand())

	return core
}
