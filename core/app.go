package core

import (
	"servon/core/internal/commands"
	"servon/core/internal/contract"
	"servon/core/internal/libs"
	"servon/core/internal/managers"
	"servon/core/internal/models"
	"servon/core/internal/utils"
	"servon/core/internal/web"
)

const DataRootFolder = "/data"
const LoggerFolder = "/logs"
const DefaultHost = "0.0.0.0"
const DefaultPort = 8080

type OSType = libs.OSType
type CommandOptions = utils.CommandOptions
type CronTask = libs.CronTask
type ValidationError = libs.ValidationError
type ValidationErrors = libs.ValidationErrors
type Task = models.Task
type SoftwareInfo = contract.SoftwareInfo
type SuperSoft = contract.SuperSoft

type App struct {
	// Libs，无状态的

	*libs.PortManager
	*libs.BasicInfoManager
	*libs.OSInfoManager
	*libs.SystemResourcesManager
	*libs.ProcessManager
	*libs.NetworkManager
	*libs.Dpkg
	*libs.CronManager
	*libs.UserLib
	*libs.TaskManager

	// Managers，有状态的或有数据的

	*managers.DeployManager
	*managers.DownloadManager
	*managers.GitManager
	*managers.AptManager
	*managers.SoftManager
	*managers.VersionManager
	*managers.CommandManager
	*managers.DataManager
	*managers.ServiceManager

	// Utils

	*utils.Printer
	*utils.CommandUtil
	*utils.FileUtil

	// WebServer

	*web.WebServerManager
}

const (
	Ubuntu  OSType = "ubuntu"
	Debian  OSType = "debian"
	CentOS  OSType = "centos"
	RedHat  OSType = "redhat"
	Unknown OSType = "unknown"
)

// New 创建Core实例
func New() *App {
	rootCmd := commands.RootCmd

	core := &App{
		CommandManager:         managers.NewCommandManager(rootCmd),
		SoftManager:            managers.DefaultSoftManager,
		DataManager:            managers.DefaultDataManager,
		Printer:                utils.DefaultPrinter,
		PortManager:            libs.DefaultPortManager,
		BasicInfoManager:       libs.DefaultBasicInfoManager,
		OSInfoManager:          libs.DefaultOSInfoManager,
		SystemResourcesManager: libs.DefaultSystemResourcesManager,
		ProcessManager:         libs.DefaultProcessManager,
		FileUtil:               utils.DefaultFileUtil,
		NetworkManager:         libs.DefaultNetworkManager,
		ServiceManager:         managers.DefaultServiceManager,
		AptManager:             managers.DefaultAptManager,
		Dpkg:                   libs.DefaultDpkg,
		CronManager:            libs.DefaultCronManager,
		VersionManager:         managers.NewVersionManager(),
		UserLib:                libs.DefaultUserManager,
		DownloadManager:        managers.NewDownloadManager(),
		GitManager:             managers.NewGitManager(),
		TaskManager:            libs.DefaultTaskManager,
		DeployManager:          managers.NewDeployManager(commands.GetDeployCommand()),
		WebServerManager: web.NewWebServerManager(
			DefaultHost,
			DefaultPort,
		),
	}

	core.AddCommand(commands.GetDeployCommand())
	core.AddCommand(commands.GetVersionCommand(core.VersionManager))
	core.AddCommand(commands.GetUpgradeCommand(core.VersionManager))
	core.AddCommand(commands.GetServiceRootCommand(core.ServiceManager))
	core.AddCommand(commands.GetUserRootCommand(core.UserLib))
	core.AddCommand(commands.GetSoftwareCommand(core.SoftManager))
	core.AddCommand(commands.GetGitRootCommand(core.GitManager))

	return core
}
