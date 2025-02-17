// Package storage 提供了 GitHub 相关数据的存储功能
// 主要负责：
// 1. 保存 webhook 事件数据
// 2. 读取历史 webhook 数据
// 3. 管理数据文件的组织结构
package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"servon/core/internal/libs/github/models"
)

// SaveWebhookPayload 保存 webhook 事件数据到指定目录
// 文件名格式：时间戳_事件ID_事件类型.json
func SaveWebhookPayload(dataDir string, eventType, eventID string, payload []byte) error {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	filename := fmt.Sprintf("%s/%d_%s_%s.json",
		dataDir,
		time.Now().Unix(),
		eventID,
		eventType,
	)

	if err := os.WriteFile(filename, payload, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// GetWebhooks 从指定目录获取所有保存的 webhook 事件数据
// 返回 WebhookPayload 数组，包含所有成功解析的事件数据
func GetWebhooks(dataDir string) ([]models.WebhookPayload, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create webhooks directory: %v", err)
	}

	files, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read webhooks directory: %v", err)
	}

	var webhooks []models.WebhookPayload
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			webhook, err := readWebhookFile(dataDir, file.Name())
			if err != nil {
				continue
			}
			webhooks = append(webhooks, webhook)
		}
	}

	return webhooks, nil
}

// readWebhookFile 读取单个 webhook 数据文件
// 解析文件名中的元数据和文件内容
func readWebhookFile(dataDir, filename string) (models.WebhookPayload, error) {
	var webhook models.WebhookPayload

	parts := strings.Split(strings.TrimSuffix(filename, ".json"), "_")
	if len(parts) != 3 {
		return webhook, fmt.Errorf("invalid filename format")
	}

	data, err := os.ReadFile(filepath.Join(dataDir, filename))
	if err != nil {
		return webhook, err
	}

	var payload interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return webhook, err
	}

	webhook.Timestamp = parseTimestamp(parts[0])
	webhook.ID = parts[1]
	webhook.Type = parts[2]
	webhook.Payload = payload

	return webhook, nil
}

// parseTimestamp 将字符串时间戳解析为 Unix 时间戳
func parseTimestamp(ts string) int64 {
	timestamp, _ := time.Parse(time.RFC3339, ts)
	return timestamp.Unix()
}
