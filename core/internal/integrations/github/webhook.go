package github

// WebhookPayload 表示存储的 webhook 事件数据
type WebhookPayload struct {
	ID        string      `json:"id"`        // 事件的唯一标识符
	Type      string      `json:"type"`      // 事件类型
	Timestamp int64       `json:"timestamp"` // 事件发生的时间戳
	Payload   interface{} `json:"payload"`   // 事件的具体内容
}
