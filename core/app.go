package core

import (
	"fmt"
	"path/filepath"
	"servon/core/internal/commands"
	"servon/core/internal/events"
	"servon/core/internal/integrations"
	"servon/core/internal/libs"
	"servon/core/internal/managers"
	"servon/core/internal/utils"
)

// App 是整个系统的核心结构，它通过事件总线协调各个功能模块之间的通信。
// 系统架构采用模块化设计，主要包含以下几个部分：
//
// 1. EventBus (事件总线)
//   - 作为系统的消息中心，负责模块间的解耦通信
//   - 支持发布/订阅模式和请求/响应模式
//   - 所有模块间的跨界通信都通过 EventBus 进行，避免直接依赖
//
// 2. Libs (无状态库)
//   - 提供基础功能支持，如端口管理、系统信息查询等
//   - 这些组件是无状态的，不需要在模块间共享数据
//   - 例如：PortManager, OSInfoManager, ProcessManager 等
//
// 3. Managers (有状态管理器)
//   - 负责特定业务领域的状态管理和功能实现
//   - 各个 Manager 之间通过 EventBus 通信，保持松耦合
//   - 例如：
//   - DeployManager: 负责部署相关的操作
//   - GitManager: 处理 Git 相关的操作
//   - ServiceManager: 管理系统服务
//
// 4. Utils (工具类)
//   - 提供通用的工具函数和辅助方法
//   - 被其他模块共同使用的基础设施
//
// 使用示例：
//
//	// 创建 App 实例
//	app := New()
//
//	// Manager 之间的通信示例
//	// 1. DeployManager 需要获取 Git 信息时：
//	app.eventBus.Request(events.Request{
//	    Type: "git.get_repo_status",
//	    Data: map[string]interface{}{"repo": "example"},
//	})
//
//	// 2. ServiceManager 在服务状态变更时通知其他模块：
//	app.eventBus.Publish(events.Event{
//	    Type: "service.status_changed",
//	    Data: map[string]interface{}{
//	        "service": "nginx",
//	        "status": "running",
//	    },
//	})
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
	*managers.WebServerManager

	// Integrations

	*integrations.GitHubIntegration

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

	deployManager, err := managers.NewDeployManager(eventBus)
	if err != nil {
		panic(fmt.Sprintf("Failed to create deploy manager: %v", err))
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
		DeployManager:          deployManager,
		WebServerManager:       managers.NewWebServerManager(DefaultHost, DefaultPort),

		// Integrations
		GitHubIntegration: integrations.NewGitHubIntegration(eventBus),

		// Utils
		DevUtil:    utils.DefaultDevUtil,
		StringUtil: utils.DefaultStringUtil,
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
