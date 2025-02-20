// Package webhook 提供了处理 GitHub webhook 事件的功能
// 主要负责：
// 1. 验证 webhook 请求的签名
// 2. 解析不同类型的 webhook 事件
// 3. 处理各种 GitHub 事件（安装、推送、PR等）
package github

import (
	"encoding/json"
	"fmt"
	"servon/core/internal/events"
	"servon/core/internal/integrations/github/models"

	"github.com/gin-gonic/gin"
)

// ProcessWebhookEvent 处理 GitHub webhook 请求
func ProcessWebhookEvent(c *gin.Context, config *models.GitHubConfig, eventBus *events.EventBus) error {
	// 首先读取 payload
	payload, err := c.GetRawData()
	if err != nil {
		return fmt.Errorf("failed to read payload: %v", err)
	}

	// 验证 webhook
	if err := validateWebhook(c, config.GitHubWebhookSecret, payload); err != nil {
		return fmt.Errorf("webhook validation failed: %v", err)
	}

	event := c.GetHeader("X-GitHub-Event")
	eventID := c.GetHeader("X-GitHub-Delivery")

	// 保存 webhook 数据
	if err := SaveWebhookPayload(WebhookDir, event, eventID, payload); err != nil {
		return fmt.Errorf("failed to save webhook payload: %v", err)
	}

	// 处理特定事件
	return handleEvent(event, payload, eventBus)
}

func validateWebhook(c *gin.Context, webhookSecret string, payload []byte) error {
	signature := c.GetHeader("X-Hub-Signature-256")
	if signature == "" {
		return fmt.Errorf("missing signature")
	}

	// TODO: 实现签名验证
	return nil
}

func handleEvent(event string, payload []byte, eventBus *events.EventBus) error {
	switch event {
	case "installation", "installation_repositories":
		return handleInstallationEvent(payload)
	case "push":
		return handlePushEvent(payload, eventBus)
	case "pull_request":
		return handlePullRequestEvent(payload)
	case "check_suite":
		return handleCheckSuiteEvent(payload)
	default:
		// 记录未处理的事件类型
		return nil
	}
}

// handleInstallationEvent 处理 GitHub App 安装相关的事件
func handleInstallationEvent(payload []byte) error {
	var event struct {
		Installation struct {
			ID      int64 `json:"id"`
			Account struct {
				Login string `json:"login"`
				ID    int64  `json:"id"`
			} `json:"account"`
			Repositories []struct {
				Name     string `json:"name"`
				FullName string `json:"full_name"`
			} `json:"repositories"`
		} `json:"installation"`
	}

	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to parse installation event: %v", err)
	}

	// 更新安装信息
	installation := &models.Installation{
		ID:           event.Installation.ID,
		AccountID:    event.Installation.Account.ID,
		AccountLogin: event.Installation.Account.Login,
		Repositories: make([]string, 0),
	}

	// 添加仓库
	for _, repo := range event.Installation.Repositories {
		installation.Repositories = append(installation.Repositories, repo.Name)
	}

	// 更新配置
	// 注意：这里需要获取 GitHubIntegration 实例来更新

	return nil
}

// handlePushEvent 处理代码推送事件
func handlePushEvent(payload []byte, eventBus *events.EventBus) error {
	// TODO: 解析 payload 并发布事件
	return eventBus.Publish(events.Event{
		Type: events.GitPush,
		Data: map[string]interface{}{
			"payload": string(payload),
		},
	})
}

// handlePullRequestEvent 处理拉取请求事件
func handlePullRequestEvent(payload []byte) error {
	// TODO: 实现PR事件处理逻辑
	return nil
}

// handleCheckSuiteEvent 处理检查套件事件
func handleCheckSuiteEvent(payload []byte) error {
	// TODO: 实现检查套件事件处理逻辑
	return nil
}
