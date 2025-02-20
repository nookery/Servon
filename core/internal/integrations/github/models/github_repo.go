package models

// GitHubRepo 表示一个GitHub仓库的基本信息
type GitHubRepo struct {
	ID          int64  `json:"id"`          // 仓库的唯一标识符
	Name        string `json:"name"`        // 仓库名称
	FullName    string `json:"full_name"`   // 仓库完整名称 (owner/name)
	Description string `json:"description"` // 仓库描述
	Private     bool   `json:"private"`     // 是否为私有仓库
	HTMLURL     string `json:"html_url"`    // 仓库的Web页面URL
}
