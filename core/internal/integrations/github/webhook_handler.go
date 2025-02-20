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

	"github.com/gin-gonic/gin"
)

// ProcessWebhookEvent 处理 GitHub webhook 请求
func (g *GitHubIntegration) ProcessWebhookEvent(c *gin.Context, config *GitHubConfig, eventBus *events.EventBus) error {
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
	return g.handleEvent(event, payload, eventBus)
}

func validateWebhook(c *gin.Context, webhookSecret string, payload []byte) error {
	signature := c.GetHeader("X-Hub-Signature-256")
	if signature == "" {
		return fmt.Errorf("missing signature")
	}

	// TODO: 实现签名验证
	return nil
}

func (g *GitHubIntegration) handleEvent(event string, payload []byte, eventBus *events.EventBus) error {
	switch event {
	case "installation", "installation_repositories":
		return g.handleInstallationEvent(payload)
	case "push":
		return g.handlePushEvent(payload, eventBus)
	case "pull_request":
		return g.handlePullRequestEvent(payload)
	case "check_suite":
		return g.handleCheckSuiteEvent(payload)
	default:
		// 记录未处理的事件类型
		return nil
	}
}

// handleInstallationEvent 处理 GitHub App 安装相关的事件
func (g *GitHubIntegration) handleInstallationEvent(payload []byte) error {
	// 首先保存原始 webhook payload
	if err := g.logger.SaveRawInstallationData(payload); err != nil {
		g.logger.LogErrorf("保存原始安装数据失败: %v", err)
		return fmt.Errorf("failed to save raw installation data: %v", err)
	}

	var event struct {
		Action       string       `json:"action"`
		Installation Installation `json:"installation"`
		Repositories []Repository `json:"repositories"`
	}

	if err := json.Unmarshal(payload, &event); err != nil {
		g.logger.LogErrorf("解析安装事件失败: %v", err)
		return fmt.Errorf("failed to parse installation event: %v", err)
	}

	// 更新安装信息
	installation := &event.Installation
	installation.Repositories = event.Repositories

	// 保存配置信息
	if err := g.logger.SaveInstallationConfig(installation); err != nil {
		g.logger.LogErrorf("保存安装配置失败: %v", err)
		return fmt.Errorf("failed to save installation config: %v", err)
	}

	// 记录安装事件数据
	installationData, err := json.MarshalIndent(installation, "", "  ")
	if err != nil {
		g.logger.LogErrorf("序列化安装数据失败: %v", err)
		return fmt.Errorf("failed to marshal installation data: %v", err)
	}

	// 保存安装数据
	if err := g.logger.SaveInstallationData(installation.ID, installationData); err != nil {
		g.logger.LogErrorf("保存安装数据失败: %v", err)
		return fmt.Errorf("failed to save installation data: %v", err)
	}

	// 使用 logger 记录安装信息
	g.logger.LogInfof("新的 GitHub App 安装: ID=%d, Account=%s",
		installation.ID,
		installation.AccountLogin,
	)

	// 更新配置
	g.config.Installations[installation.ID] = installation

	return nil
}

// handlePushEvent 处理代码推送事件
func (g *GitHubIntegration) handlePushEvent(payload []byte, eventBus *events.EventBus) error {
	// TODO: 解析 payload 并发布事件
	return eventBus.Publish(events.Event{
		Type: events.GitPush,
		Data: map[string]interface{}{
			"payload": string(payload),
		},
	})
}

// handlePullRequestEvent 处理拉取请求事件
func (g *GitHubIntegration) handlePullRequestEvent(payload []byte) error {
	// TODO: 实现PR事件处理逻辑
	return nil
}

// handleCheckSuiteEvent 处理检查套件事件
func (g *GitHubIntegration) handleCheckSuiteEvent(payload []byte) error {
	// TODO: 实现检查套件事件处理逻辑
	return nil
}
