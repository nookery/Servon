package providers

import (
	"servon/core/internal/events"
	"servon/core/internal/managers"
)

type ManagerProvider struct {
	*managers.FullManager
}

func NewManagerProvider(eventBus *events.EventBus, manager *managers.FullManager) *ManagerProvider {
	core := &ManagerProvider{
		FullManager: manager,
	}

	return core
}
