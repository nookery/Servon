package managers

import (
	"fmt"
	"servon/core/internal/events"
	"servon/core/internal/managers/github"
)

type FullManager struct {
	*CronManager
	*DeployManager
	*DownloadManager
	*GitManager
	*AptManager
	*SoftManager
	*VersionManager
	*DataManager
	*ServiceManager
	*FileManager
	*OSInfoManager
	*SystemResourcesManager
	*BasicInfoManager
	*NetworkManager
	*PortManager
	*UserManager
	*TaskManager
	*ProcessManager
	*DpkgManager
	*LogManager
	*github.GitHubIntegration
}

func NewManager(eventBus *events.EventBus) *FullManager {
	dataManager := DefaultDataManager
	githubIntegration := github.NewGitHubIntegration(eventBus)
	deployManager, err := NewDeployManager(eventBus, githubIntegration, dataManager.GetLogsRootFolder())
	if err != nil {
		panic(fmt.Sprintf("Failed to create deploy manager: %v", err))
	}

	core := &FullManager{
		CronManager:            DefaultCronManager,
		SoftManager:            DefaultSoftManager,
		DataManager:            dataManager,
		ServiceManager:         DefaultServiceManager,
		AptManager:             DefaultAptManager,
		VersionManager:         NewVersionManager(),
		DownloadManager:        NewDownloadManager(),
		GitManager:             NewGitManager(),
		DeployManager:          deployManager,
		FileManager:            NewFileManager(),
		OSInfoManager:          NewOSInfoManager(),
		SystemResourcesManager: NewSystemResourcesManager(),
		BasicInfoManager:       NewBasicInfoManager(),
		NetworkManager:         NewNetworkManager(),
		PortManager:            NewPortManager(),
		TaskManager:            DefaultTaskManager,
		UserManager:            NewUserManager(),
		ProcessManager:         DefaultProcessManager,
		GitHubIntegration:      githubIntegration,
		LogManager:             DefaultLogManager,
	}

	return core
}
