// Package models 定义了 GitHub 相关功能使用的数据结构
package github

// Installation 表示 GitHub App 的安装信息
type Installation struct {
	ID           int64    `json:"id"`            // 安装的唯一标识符
	AccountID    int64    `json:"account_id"`    // 安装账户的ID
	AccountLogin string   `json:"account_login"` // 安装账户的登录名
	Repositories []string `json:"repositories"`  // 安装的仓库列表
}
