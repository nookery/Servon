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

// GetInstallationConfig 从存储中读取所有安装配置信息
func GetInstallationConfig() (map[int64]*Installation, error) {
	// 确保目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 读取配置目录下的所有文件
	files, err := os.ReadDir(configDir)
	if err != nil {
		return nil, fmt.Errorf("读取配置目录失败: %w", err)
	}

	installations := make(map[int64]*Installation)

	// 遍历所有配置文件
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), "installation_") || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(configDir, file.Name()))
		if err != nil {
			printer.PrintErrorf("读取配置文件失败 %s: %v", file.Name(), err)
			continue
		}

		var config InstallationConfig
		if err := json.Unmarshal(data, &config); err != nil {
			printer.PrintErrorf("解析配置文件失败 %s: %v", file.Name(), err)
			continue
		}

		// 转换为 Installation 对象
		installation := &Installation{
			ID:           config.InstallationID,
			AccountID:    config.AccountID,
			AccountLogin: config.AccountLogin,
			AccountType:  config.AccountType,
			AppID:        config.AppID,
			Permissions:  config.Permissions,
			Events:       config.Events,
			Repositories: config.Repositories,
			CreatedAt:    config.CreatedAt,
		}

		installations[config.InstallationID] = installation
	}

	return installations, nil
}

// SaveInstallationData 保存安装数据到指定目录
func SaveInstallationData(installationID int64, data []byte) error {
	// 确保目录存在
	if err := os.MkdirAll(installationDir, 0755); err != nil {
		return fmt.Errorf("创建安装数据目录失败: %v", err)
	}

	filename := filepath.Join(installationDir, fmt.Sprintf("%d.json", installationID))
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("写入安装数据失败: %v", err)
	}

	return nil
}

// SaveRawInstallationData 保存原始安装数据，使用时间戳作为文件名
func SaveRawInstallationData(payload []byte) error {
	// 确保目录存在
	if err := os.MkdirAll(installationDir, 0755); err != nil {
		return fmt.Errorf("创建安装数据目录失败: %v", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(installationDir, fmt.Sprintf("raw_%s.json", timestamp))

	if err := os.WriteFile(filename, payload, 0644); err != nil {
		return fmt.Errorf("写入原始安装数据失败: %v", err)
	}

	return nil
}

// SaveInstallationConfig 保存安装配置到指定目录
func SaveInstallationConfig(installation *Installation) error {
	config := InstallationConfig{
		InstallationID: installation.ID,
		AccountID:      installation.AccountID,
		AccountLogin:   installation.AccountLogin,
		AccountType:    installation.AccountType,
		AppID:          installation.AppID,
		Permissions:    installation.Permissions,
		Events:         installation.Events,
		Repositories:   installation.Repositories,
		CreatedAt:      installation.CreatedAt,
	}

	// 确保目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// 生成配置文件路径
	configPath := filepath.Join(configDir, fmt.Sprintf("installation_%d.json", installation.ID))

	// 序列化并保存配置
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// SaveAppConfig 保存 GitHub App 配置到磁盘
func SaveAppConfig(config *GitHubConfig) error {
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	config.UpdatedAt = time.Now().Format(time.RFC3339)
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	configPath := filepath.Join(configDir, "app_config.json")
	if err := os.WriteFile(configPath, data, 0600); err != nil { // 使用更严格的权限
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}

// LoadAppConfig 从磁盘加载 GitHub App 配置
func LoadAppConfig() (*GitHubConfig, error) {
	configPath := filepath.Join(configDir, "app_config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // 配置文件不存在返回 nil
		}
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config GitHubConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}
