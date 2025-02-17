// Package models 定义了 GitHub 相关功能使用的数据结构
// 主要包含：
// 1. GitHub App 安装信息
// 2. GitHub App 配置信息
// 3. GitHub App manifest 配置
// 4. Webhook 事件数据结构
package models

// Installation 表示 GitHub App 的安装信息
type Installation struct {
	ID           int64    `json:"id"`            // 安装的唯一标识符
	AccountID    int64    `json:"account_id"`    // 安装账户的ID
	AccountLogin string   `json:"account_login"` // 安装账户的登录名
	Repositories []string `json:"repositories"`  // 安装的仓库列表
}

// GitHubConfig 表示 GitHub App 的配置信息
type GitHubConfig struct {
	GitHubAppID         int64                   `json:"github_app_id"`          // GitHub App 的ID
	GitHubAppPrivateKey string                  `json:"github_app_private_key"` // GitHub App 的私钥
	GitHubWebhookSecret string                  `json:"github_webhook_secret"`  // Webhook 的密钥
	Installations       map[int64]*Installation `json:"installations"`          // 所有安装实例的映射表
}

// GitHubManifest 表示 GitHub App 的 manifest 配置
// 用于创建新的 GitHub App
type GitHubManifest struct {
	Name           string `json:"name"` // 应用名称
	URL            string `json:"url"`  // 应用主页URL
	HookAttributes struct {
		URL    string `json:"url"`    // Webhook 接收URL
		Active bool   `json:"active"` // Webhook 是否激活
	} `json:"hook_attributes"`
	RedirectURL        string            `json:"redirect_url"`        // 授权后的重定向URL
	CallbackURLs       []string          `json:"callback_urls"`       // 回调URL列表
	Description        string            `json:"description"`         // 应用描述
	Public             bool              `json:"public"`              // 是否公开应用
	DefaultEvents      []string          `json:"default_events"`      // 默认订阅的事件列表
	DefaultPermissions map[string]string `json:"default_permissions"` // 默认请求的权限列表
}

// WebhookPayload 表示存储的 webhook 事件数据
type WebhookPayload struct {
	ID        string      `json:"id"`        // 事件的唯一标识符
	Type      string      `json:"type"`      // 事件类型
	Timestamp int64       `json:"timestamp"` // 事件发生的时间戳
	Payload   interface{} `json:"payload"`   // 事件的具体内容
}
