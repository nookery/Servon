// Package models 定义了 GitHub 相关功能使用的数据结构
package github

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
