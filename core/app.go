package core

import (
	"fmt"
	"path/filepath"
	"servon/core/internal/events"
	"servon/core/internal/managers"
	"servon/core/internal/providers"
)

type App struct {
	eventBus *events.EventBus

	*providers.WebProvider
	*providers.ManagerProvider
	*providers.CommandProvider
	*providers.UtilProvider
}

// New 创建App实例
func New() *App {
	eventBus, err := events.NewEventBus(filepath.Join(DataRootFolder, "events"))
	if err != nil {
		panic(fmt.Sprintf("Failed to create event bus: %v", err))
	}

	manager := managers.NewManager(eventBus)
	webProvider := providers.NewWebProvider(manager, DefaultHost, DefaultPort)

	app := &App{
		eventBus:        eventBus,
		WebProvider:     webProvider,
		ManagerProvider: providers.NewManagerProvider(eventBus, manager),
		CommandProvider: providers.NewCommandProvider(manager, webProvider.Server),
		UtilProvider:    providers.NewUtilProvider(),
	}

	app.Success("App 初始化完成")

	return app
}
