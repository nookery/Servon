// Package managers 提供了各种管理器的实现
// github_manager.go 负责管理 GitHub App 的整个生命周期
// 主要功能：
// 1. GitHub App 的安装和配置
// 2. Webhook 事件的处理和存储
// 3. 配置信息的管理和持久化
package managers

import (
	"fmt"
	"sync"

	"servon/core/internal/events"
	"servon/core/internal/libs/github/config"
	"servon/core/internal/libs/github/models"
	"servon/core/internal/libs/github/storage"
	"servon/core/internal/libs/github/webhook"

	"github.com/gin-gonic/gin"
)

// GitHubManager 管理 GitHub App 的所有操作
// 使用单例模式确保全局只有一个实例
type GitHubManager struct {
	config   *models.GitHubConfig // GitHub App 的配置信息
	mu       sync.RWMutex         // 用于保护配置访问的互斥锁
	eventBus *events.EventBus     // 添加事件总线
}

var (
	instance *GitHubManager // 单例实例
	once     sync.Once      // 确保单例只初始化一次
)

// dataDir 指定 webhook 数据存储的目录
const dataDir = "/data/github" // TODO: 从配置中读取

// GetGitHubManager 返回 GitHubManager 的单例实例
// 确保全局只有一个 GitHubManager 实例在运行
func GetGitHubManager(eventBus *events.EventBus) *GitHubManager {
	once.Do(func() {
		instance = &GitHubManager{
			config:   &models.GitHubConfig{Installations: make(map[int64]*models.Installation)},
			eventBus: eventBus,
		}
	})
	return instance
}

// HandleSetup 处理 GitHub App 的安装设置
// 生成必要的 manifest 配置并返回安装页面
func (m *GitHubManager) HandleSetup(c *gin.Context) (string, error) {
	c.Set("base_url", "http://43.142.208.212:9754") // TODO: 从配置中读取
	return config.GenerateManifest(c)
}

// HandleCallback 处理 GitHub App 安装后的回调
// 保存配置信息并返回安装URL
func (m *GitHubManager) HandleCallback(c *gin.Context) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	result, err := config.ProcessCallback(c)
	if err != nil {
		return "", err
	}

	m.config.GitHubAppID = result.ID
	m.config.GitHubAppPrivateKey = result.PEM

	return result.GetInstallURL(), nil
}

// HandleWebhook 处理接收到的 GitHub webhook 事件
func (m *GitHubManager) HandleWebhook(c *gin.Context) error {
	printer.PrintInfof("HandleWebhook")
	if err := webhook.HandleWebhook(c, m.config.GitHubWebhookSecret); err != nil {
		printer.PrintErrorf("HandleWebhook error: %v", err)
		return err
	}

	payload, err := c.GetRawData()
	if err != nil {
		return fmt.Errorf("failed to read payload: %v", err)
	}

	event := c.GetHeader("X-GitHub-Event")
	eventID := c.GetHeader("X-GitHub-Delivery")

	// 优先保存 webhook 数据
	if err := storage.SaveWebhookPayload(dataDir, event, eventID, payload); err != nil {
		printer.PrintErrorf("failed to save webhook payload: %v", err)
		return fmt.Errorf("failed to save webhook payload: %v", err)
	}

	// 数据保存成功后，处理特定事件
	if event == "push" {
		// 这里使用样本数据，实际应该从payload中解析
		sampleDeployData := map[string]interface{}{
			"repository": "https://github.com/yourusername/yourrepo",
			"branch":     "main",
			"commit":     "abc123",
		}

		// 发布部署事件
		if err := m.eventBus.Publish(events.Event{
			Type: events.GitPush,
			Data: sampleDeployData,
		}); err != nil {
			printer.PrintErrorf("failed to publish deploy event: %v", err)
			// 注意：这里的错误不会影响webhook数据的保存，因为数据已经保存成功
			return err
		}
	}

	return nil
}

// GetStoredWebhooks 获取所有存储的 webhook 事件数据
func (m *GitHubManager) GetStoredWebhooks() ([]models.WebhookPayload, error) {
	return storage.GetWebhooks(dataDir)
}

// GetConfig 返回当前的 GitHub 配置信息
// 使用读锁保护配置访问
func (m *GitHubManager) GetConfig() *models.GitHubConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}
