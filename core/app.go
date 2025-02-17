package core

import (
	"fmt"
	"path/filepath"
	"servon/core/internal/commands"
	"servon/core/internal/events"
	"servon/core/internal/libs"
	"servon/core/internal/managers"
	"servon/core/internal/utils"
)

type App struct {
	eventBus *events.EventBus

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
	*managers.GitHubManager
	*managers.WebServerManager

	// Utils

	*utils.Printer
	*utils.CommandUtil
	*utils.FileUtil
	*utils.DevUtil
	*utils.StringUtil
}

// New 创建Core实例
func New() *App {
	rootCmd := commands.RootCmd
	eventBus, err := events.NewEventBus(filepath.Join(DataRootFolder, "events"))
	if err != nil {
		panic(fmt.Sprintf("Failed to create event bus: %v", err))
	}

	core := &App{
		eventBus:               eventBus,
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
		GitHubManager:          managers.GetGitHubManager(eventBus),
		DeployManager:          managers.NewDeployManager(eventBus),
		WebServerManager:       managers.NewWebServerManager(DefaultHost, DefaultPort),
		DevUtil:                utils.DefaultDevUtil,
		StringUtil:             utils.DefaultStringUtil,
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
