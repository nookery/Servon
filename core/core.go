package core

import (
	"servon/core/libs"
)

// 调用关系 Core -> Core API -> Libs
// 或 Core -> Libs

const DataRootFolder = "/data"
const LoggerFolder = "/logs"

type Core struct {
	*libs.CommandManager
	*libs.DataManager
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
	*libs.ShellManager
	*libs.UserManager
	*libs.EnvManager
}

type OSType = libs.OSType
type CommandOptions = libs.CommandOptions
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
		CommandManager:         libs.DefaultCommandManager,
		SoftManager:            libs.DefaultSoftManager,
		DataManager:            libs.DefaultDataManager,
		DeployManager:          libs.DefaultDeployManager,
		Printer:                libs.DefaultPrinter,
		PortManager:            libs.DefaultPortManager,
		BasicInfoManager:       libs.DefaultBasicInfoManager,
		OSInfoManager:          libs.DefaultOSInfoManager,
		SystemResourcesManager: libs.DefaultSystemResourcesManager,
		ProcessManager:         libs.DefaultProcessManager,
		FilesManager:           libs.DefaultFilesManager,
		NetworkManager:         libs.DefaultNetworkManager,
		ServiceManager:         libs.DefaultServiceManager,
		AptManager:             libs.DefaultAptManager,
		Dpkg:                   libs.DefaultDpkg,
		CronManager:            libs.DefaultCronManager,
		VersionManager:         libs.DefaultVersionManager,
		ShellManager:           libs.DefaultShellManager,
		UserManager:            libs.DefaultUserManager,
		EnvManager:             libs.DefaultEnvManager,
	}

	core.LoadEnv()

	core.AddCommand(core.GetDeployCommand())
	core.AddCommand(core.GetVersionCommand())
	core.AddCommand(core.GetSoftwareCommand())
	core.AddCommand(core.GetUserRootCommand())

	return core
}
