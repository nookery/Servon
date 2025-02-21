package github

import (
	"time"

	"github.com/fatih/color"
)

// GitHubRepo 表示一个GitHub仓库的基本信息
type GitHubRepo struct {
	ID          int64  `json:"id"`          // 仓库的唯一标识符
	Name        string `json:"name"`        // 仓库名称
	FullName    string `json:"full_name"`   // 仓库完整名称 (owner/name)
	Description string `json:"description"` // 仓库描述
	Private     bool   `json:"private"`     // 是否为私有仓库
	HTMLURL     string `json:"html_url"`    // 仓库的Web页面URL
}

// InstallationConfig 表示GitHub App安装的配置信息
type InstallationConfig struct {
	InstallationID int64        `json:"installation_id"` // 安装的唯一标识符
	AccountID      int64        `json:"account_id"`      // 账户ID
	AccountLogin   string       `json:"account_login"`   // 账户登录名
	AccountType    string       `json:"account_type"`    // 账户类型
	AppID          int64        `json:"app_id"`          // GitHub App ID
	Permissions    Permissions  `json:"permissions"`     // 权限配置
	Events         []string     `json:"events"`          // 订阅的事件
	Repositories   []Repository `json:"repositories"`    // 仓库列表
	CreatedAt      string       `json:"created_at"`      // 创建时间
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

// GitHubConfig 表示 GitHub App 的配置信息
type GitHubConfig struct {
	GitHubAppID         int64                   `json:"github_app_id"`          // GitHub App 的ID
	GitHubAppPrivateKey string                  `json:"github_app_private_key"` // GitHub App 的私钥
	GitHubWebhookSecret string                  `json:"github_webhook_secret"`  // Webhook 的密钥
	Installations       map[int64]*Installation `json:"installations"`          // 所有安装实例的映射表
}

// AppConfig 存储 GitHub App 的基本配置
type AppConfig struct {
	AppID      int64  `json:"app_id"`      // GitHub App 的ID
	PrivateKey string `json:"private_key"` // GitHub App 的私钥
	WebhookKey string `json:"webhook_key"` // Webhook 密钥
	UpdatedAt  string `json:"updated_at"`  // 最后更新时间
}

// Installation 表示 GitHub App 的安装信息
type Installation struct {
	ID                  int64        `json:"id"`                        // 安装的唯一标识符
	ClientID            string       `json:"client_id"`                 // 客户端ID
	AccountID           int64        `json:"account_id"`                // 安装账户的ID
	AccountLogin        string       `json:"account_login"`             // 安装账户的登录名
	AccountType         string       `json:"account_type"`              // 账户类型(User/Organization)
	AccountAvatarURL    string       `json:"account_avatar_url"`        // 账户头像URL
	RepositorySelection string       `json:"repository_selection"`      // 仓库选择模式(all/selected)
	AccessTokensURL     string       `json:"access_tokens_url"`         // 访问令牌URL
	RepositoriesURL     string       `json:"repositories_url"`          // 仓库API URL
	HTMLURL             string       `json:"html_url"`                  // 安装页面URL
	AppID               int64        `json:"app_id"`                    // GitHub App的ID
	AppSlug             string       `json:"app_slug"`                  // GitHub App的标识符
	TargetID            int64        `json:"target_id"`                 // 目标账户ID
	TargetType          string       `json:"target_type"`               // 目标类型
	SingleFileName      *string      `json:"single_file_name"`          // 单文件名称(如果适用)
	HasMultipleFiles    bool         `json:"has_multiple_single_files"` // 是否有多个单文件
	SingleFilePaths     []string     `json:"single_file_paths"`         // 单文件路径列表
	Repositories        []Repository `json:"repositories"`              // 安装的仓库列表
	Permissions         Permissions  `json:"permissions"`               // 安装权限
	Events              []string     `json:"events"`                    // 订阅的事件列表
	CreatedAt           string       `json:"created_at"`                // 创建时间
	UpdatedAt           string       `json:"updated_at"`                // 更新时间
	SuspendedBy         *string      `json:"suspended_by"`              // 暂停者(如果被暂停)
	SuspendedAt         *string      `json:"suspended_at"`              // 暂停时间(如果被暂停)
	Account             struct {
		Login     string `json:"login"`      // 账户登录名
		ID        int64  `json:"id"`         // 账户ID
		NodeID    string `json:"node_id"`    // GraphQL节点ID
		AvatarURL string `json:"avatar_url"` // 头像URL
		Type      string `json:"type"`       // 账户类型
	} `json:"account"`
}

// Repository 表示仓库信息
type Repository struct {
	ID       int64  `json:"id"`        // 仓库ID
	NodeID   string `json:"node_id"`   // GraphQL节点ID
	Name     string `json:"name"`      // 仓库名称
	FullName string `json:"full_name"` // 完整仓库名称
	Private  bool   `json:"private"`   // 是否为私有仓库
}

// Permissions 表示安装的权限配置
type Permissions struct {
	Checks   string `json:"checks"`   // 检查权限级别
	Issues   string `json:"issues"`   // Issue权限级别
	Metadata string `json:"metadata"` // 元数据权限级别
}

// WebhookPayload 表示存储的 webhook 事件数据
type WebhookPayload struct {
	ID        string      `json:"id"`        // 事件的唯一标识符
	Type      string      `json:"type"`      // 事件类型
	Timestamp int64       `json:"timestamp"` // 事件发生的时间戳
	Payload   interface{} `json:"payload"`   // 事件的具体内容
}

// TokenCache 用于缓存安装令牌
type TokenCache struct {
	Token     string    `json:"token"`      // GitHub安装令牌
	ExpiresAt time.Time `json:"expires_at"` // 令牌过期时间
}

// LogType 定义日志类型及其属性
type LogType struct {
	Name   string
	Color  *color.Color
	Symbol string
}

// 定义所有日志类型
var (
	LogTypeInfo = LogType{
		Name:   "info",
		Color:  color.New(color.FgCyan),
		Symbol: "🍋",
	}
	LogTypeError = LogType{
		Name:   "error",
		Color:  color.New(color.FgRed),
		Symbol: "❌",
	}
	LogTypeWarn = LogType{
		Name:   "warn",
		Color:  color.New(color.FgYellow),
		Symbol: "🚨",
	}
	LogTypeSuccess = LogType{
		Name:   "success",
		Color:  color.New(color.FgGreen),
		Symbol: "✅",
	}
	LogTypeDebug = LogType{
		Name:   "debug",
		Color:  color.New(color.FgBlue),
		Symbol: "🔍",
	}
)
