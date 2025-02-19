// Package storage 提供了 GitHub 相关数据的存储功能
// 主要负责：
// 1. 保存 webhook 事件数据
// 2. 读取历史 webhook 数据
// 3. 管理数据文件的组织结构
package github

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"servon/core/internal/utils"
)

var printer = utils.DefaultPrinter

const (
	WebhookDir = "/data/github/webhook"
)

// SaveWebhookPayload 保存 webhook 事件数据到指定目录
// 文件名格式：时间戳_事件ID_事件类型.json
func SaveWebhookPayload(dataDir string, eventType, eventID string, payload []byte) error {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	filename := fmt.Sprintf("%d_%s_%s.json",
		time.Now().Unix(),
		eventID,
		eventType,
	)

	return os.WriteFile(filepath.Join(dataDir, filename), payload, 0644)
}

// GetWebhooks 从指定目录获取所有保存的 webhook 事件数据
// 返回 WebhookPayload 数组，包含所有成功解析的事件数据
func GetWebhooks(dataDir string) ([]WebhookPayload, error) {
	printer.PrintInfof("GetWebhooks: %s", dataDir)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		printer.PrintErrorf("failed to create webhooks directory: %v", err)
		return nil, fmt.Errorf("failed to create webhooks directory: %v", err)
	}

	files, err := os.ReadDir(dataDir)
	if err != nil {
		printer.PrintErrorf("failed to read webhooks directory: %v", err)
		return nil, fmt.Errorf("failed to read webhooks directory: %v", err)
	}

	var webhooks []WebhookPayload
	for _, file := range files {
		printer.PrintInfof("GetWebhooks: %s", file.Name())
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dataDir, file.Name()))
		if err != nil {
			continue
		}

		var webhook WebhookPayload
		if err := json.Unmarshal(data, &webhook); err != nil {
			continue
		}

		webhooks = append(webhooks, webhook)
	}

	printer.PrintInfof("GetWebhooks: %d", len(webhooks))
	return webhooks, nil
}

// readWebhookFile 读取单个 webhook 数据文件
// 解析文件名中的元数据和文件内容
func readWebhookFile(dataDir, filename string) (WebhookPayload, error) {
	var webhook WebhookPayload

	// 移除 .json 后缀
	basename := strings.TrimSuffix(filename, ".json")

	// 找到第一个和第二个下划线的位置
	firstUnderscore := strings.Index(basename, "_")
	secondUnderscore := strings.Index(basename[firstUnderscore+1:], "_") + firstUnderscore + 1

	if firstUnderscore == -1 || secondUnderscore == -1 {
		printer.PrintErrorf("readWebhookFile: invalid filename format: %s", filename)
		return webhook, fmt.Errorf("invalid filename format")
	}

	// 正确提取三个部分
	timestamp := basename[:firstUnderscore]
	eventID := basename[firstUnderscore+1 : secondUnderscore]
	eventType := basename[secondUnderscore+1:]

	data, err := os.ReadFile(filepath.Join(dataDir, filename))
	if err != nil {
		printer.PrintErrorf("readWebhookFile: failed to read file: %v", err)
		return webhook, err
	}

	var payload interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		printer.PrintErrorf("readWebhookFile: failed to unmarshal file: %v", err)
		return webhook, err
	}

	webhook.Timestamp = parseTimestamp(timestamp)
	webhook.ID = eventID
	webhook.Type = eventType
	webhook.Payload = payload

	return webhook, nil
}

// parseTimestamp 将字符串时间戳解析为 Unix 时间戳
func parseTimestamp(ts string) int64 {
	timestamp, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		printer.PrintErrorf("parseTimestamp: failed to parse timestamp: %v", err)
		return 0
	}
	return timestamp
}
