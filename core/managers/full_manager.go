package managers

import (
	"fmt"
	"servon/components/events"
	"servon/components/github"
	"servon/components/user"
)

type FullManager struct {
	*CronManager
	*DeployManager
	*DownloadManager
	*GitManager
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
	*user.UserManager
	*TaskManager
	*ProcessManager
	*LogManager
	*ProjectManager
	*github.GitHubIntegration
}

func NewManager(eventBus events.IEventBus) *FullManager {
	dataManager := DefaultDataManager
	githubIntegration := github.NewGitHubIntegration(eventBus)
	softManager := NewSoftManager(dataManager.GetLogsRootFolder())
	gitManager := NewGitManager(softManager)
	downloadManager := NewDownloadManager(softManager)

	deployManager, err := NewDeployManager(
		eventBus,
		githubIntegration,
		dataManager.GetLogsRootFolder(),
		dataManager.GetTempRootFolder(),
		dataManager.GetProjectsFolder(),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create deploy manager: %v", err))
	}

	core := &FullManager{
		CronManager:            DefaultCronManager,
		SoftManager:            softManager,
		DataManager:            dataManager,
		ServiceManager:         DefaultServiceManager,
		VersionManager:         NewVersionManager(),
		DownloadManager:        downloadManager,
		GitManager:             gitManager,
		DeployManager:          deployManager,
		FileManager:            NewFileManager(),
		OSInfoManager:          NewOSInfoManager(),
		SystemResourcesManager: NewSystemResourcesManager(),
		BasicInfoManager:       NewBasicInfoManager(),
		NetworkManager:         NewNetworkManager(),
		PortManager:            NewPortManager(),
		TaskManager:            DefaultTaskManager,
		UserManager:            user.NewUserManager(),
		ProcessManager:         DefaultProcessManager,
		GitHubIntegration:      githubIntegration,
		LogManager:             DefaultLogManager,
		ProjectManager:         NewTopologyManager(softManager),
	}

	return core
}
