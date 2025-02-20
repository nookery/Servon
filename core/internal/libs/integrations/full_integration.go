package integrations

import "servon/core/internal/events"

type FullIntegration struct {
	*GitHubIntegration
}

func NewFullIntegration(eventBus *events.EventBus) *FullIntegration {
	return &FullIntegration{
		GitHubIntegration: NewGitHubIntegration(eventBus),
	}
}
