package models

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
