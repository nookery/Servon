package core

import (
	"path/filepath"

	"servon/components"
	"servon/components/events"
	"servon/core/internal/managers"
	"servon/core/internal/providers"
	"servon/core/internal/utils"
)

type App struct {
	eventBus events.IEventBus

	*providers.WebProvider
	*providers.ManagerProvider
	*providers.CommandProvider
	*providers.UtilProvider

	SoftwareLogger *utils.LogUtil
	AppLogger      *utils.LogUtil
}

// New 创建App实例
func New() *App {
	eventBus := components.EventBus

	manager := managers.NewManager(eventBus)
	webProvider := providers.NewWebProvider(manager, DefaultHost, DefaultPort)

	app := &App{
		eventBus:        eventBus,
		WebProvider:     webProvider,
		ManagerProvider: providers.NewManagerProvider(eventBus, manager),
		CommandProvider: providers.NewCommandProvider(manager, webProvider.Server),
		UtilProvider:    providers.NewUtilProvider(),
		SoftwareLogger:  manager.SoftManager.LogUtil,
		AppLogger:       utils.NewLogUtil(filepath.Join(DataRootFolder, "logs")),
	}

	app.AppLogger.Success("App 初始化完成")

	return app
}
