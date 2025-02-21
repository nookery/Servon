package models

import "time"

// DeployLog 表示一条部署日志记录
// 包含部署的ID、时间戳、内容、仓库信息、状态和消息等信息
type DeployLog struct {
	ID        string    `json:"id"`        // 部署日志的唯一标识符
	Timestamp time.Time `json:"timestamp"` // 部署时间
	Content   string    `json:"content"`   // 部署的详细内容
	Repo      string    `json:"repo"`      // 关联的代码仓库
	Status    string    `json:"status"`    // 部署状态（如：成功、失败）
	Message   string    `json:"message"`   // 部署相关的消息或错误信息
	// 根据实际需求添加其他字段
}
