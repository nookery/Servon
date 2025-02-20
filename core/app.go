package core

import (
	"fmt"
	"path/filepath"
	"servon/core/internal/events"
	"servon/core/internal/libs/integrations"
	"servon/core/internal/libs/managers"
	"servon/core/internal/providers"
)

type App struct {
	eventBus *events.EventBus

	*providers.WebProvider
	*providers.ManagerProvider
	*providers.IntegrationProvider
	*providers.CommandProvider
	*providers.UtilProvider
}

// New 创建Core实例
func New() *App {
	eventBus, err := events.NewEventBus(filepath.Join(DataRootFolder, "events"))
	if err != nil {
		panic(fmt.Sprintf("Failed to create event bus: %v", err))
	}

	manager := managers.NewManager(eventBus)
	integrations := integrations.NewFullIntegration(eventBus)
	webProvider := providers.NewWebProvider(manager, integrations, DefaultHost, DefaultPort)

	core := &App{
		eventBus:            eventBus,
		WebProvider:         webProvider,
		ManagerProvider:     providers.NewManagerProvider(eventBus),
		IntegrationProvider: providers.NewIntegrationProvider(eventBus),
		CommandProvider:     providers.NewCommandProvider(manager, webProvider.Server),
		UtilProvider:        providers.NewUtilProvider(),
	}

	return core
}
