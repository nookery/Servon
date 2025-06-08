package providers

import (
	"servon/components/events"
	"servon/core/managers"
)

type ManagerProvider struct {
	*managers.FullManager
}

func NewManagerProvider(eventBus events.IEventBus, manager *managers.FullManager) *ManagerProvider {
	core := &ManagerProvider{
		FullManager: manager,
	}

	return core
}
