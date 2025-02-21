package github

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
