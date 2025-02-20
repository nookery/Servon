package integrations

import (
	"servon/core/internal/events"
	"servon/core/internal/libs/github"
)

type FullIntegration struct {
	*github.GitHubIntegration
}

func NewFullIntegration(eventBus *events.EventBus) *FullIntegration {
	return &FullIntegration{
		GitHubIntegration: github.NewGitHubIntegration(eventBus),
	}
}
