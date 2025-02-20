package providers

import (
	"servon/core/internal/events"
	"servon/core/internal/libs/integrations"
)

type IntegrationProvider struct {
	*integrations.FullIntegration
}

func NewIntegrationProvider(eventBus *events.EventBus) *IntegrationProvider {
	return &IntegrationProvider{
		FullIntegration: integrations.NewFullIntegration(eventBus),
	}
}
