package models

// Installation represents a GitHub App installation
type Installation struct {
	ID           int64    `json:"id"`
	AccountID    int64    `json:"account_id"`
	AccountLogin string   `json:"account_login"`
	Repositories []string `json:"repositories"`
}

// GitHubConfig represents the GitHub App configuration
type GitHubConfig struct {
	GitHubAppID         int64                   `json:"github_app_id"`
	GitHubAppPrivateKey string                  `json:"github_app_private_key"`
	GitHubWebhookSecret string                  `json:"github_webhook_secret"`
	Installations       map[int64]*Installation `json:"installations"`
}

// GitHubManifest represents the GitHub App manifest configuration
type GitHubManifest struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	HookAttributes struct {
		URL    string `json:"url"`
		Active bool   `json:"active"`
	} `json:"hook_attributes"`
	RedirectURL        string            `json:"redirect_url"`
	CallbackURLs       []string          `json:"callback_urls"`
	Description        string            `json:"description"`
	Public             bool              `json:"public"`
	DefaultEvents      []string          `json:"default_events"`
	DefaultPermissions map[string]string `json:"default_permissions"`
}

// WebhookPayload represents a stored webhook event
type WebhookPayload struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Timestamp int64       `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}
