package providers

import (
	"servon/core/internal/events"
	"servon/core/internal/managers"
)

type ManagerProvider struct {
	*managers.FullManager
}

func NewManagerProvider(eventBus *events.EventBus) *ManagerProvider {
	core := &ManagerProvider{
		FullManager: managers.NewManager(eventBus),
	}

	return core
}
