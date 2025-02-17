package models

type DeployLog struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	// 根据实际需求添加其他字段
}
