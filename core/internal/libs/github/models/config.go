package models

// GitHubConfig 表示 GitHub App 的配置信息
type GitHubConfig struct {
	GitHubAppID         int64                   `json:"github_app_id"`          // GitHub App 的ID
	GitHubAppPrivateKey string                  `json:"github_app_private_key"` // GitHub App 的私钥
	GitHubWebhookSecret string                  `json:"github_webhook_secret"`  // Webhook 的密钥
	Installations       map[int64]*Installation `json:"installations"`          // 所有安装实例的映射表
}
