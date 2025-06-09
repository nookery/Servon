package core

import (
	"path/filepath"

	"servon/components"
	"servon/components/events"
	"servon/components/logger"
	"servon/core/managers"
	"servon/core/providers"
)

type App struct {
	eventBus events.IEventBus

	*providers.WebProvider
	*providers.ManagerProvider
	*providers.CommandProvider
	*providers.UtilProvider

	AppLogger *logger.LogUtil
}

// New 创建App实例
func New() *App {
	eventBus := components.EventBus

	manager := managers.NewManager(eventBus)
	webProvider := providers.NewWebProvider(manager)

	app := &App{
		eventBus:        eventBus,
		WebProvider:     webProvider,
		ManagerProvider: providers.NewManagerProvider(eventBus, manager),
		CommandProvider: providers.NewCommandProvider(manager, webProvider.Server),
		UtilProvider:    providers.NewUtilProvider(),
		AppLogger:       logger.NewLogUtil(filepath.Join(DataRootFolder, "logs")),
	}

	return app
}
