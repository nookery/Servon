// Package integrations 提供了各种第三方服务的集成实现
package integrations

import (
	"context"
	"servon/core/internal/events"
	"servon/core/internal/libs/github/config"
	"servon/core/internal/libs/github/logger"
	githubModels "servon/core/internal/libs/github/models"
	"servon/core/internal/libs/github/storage"
	"servon/core/internal/libs/github/webhook"
	"servon/core/internal/models"
	"servon/core/internal/utils"
	"sync"

	"github.com/gin-gonic/gin"
)

var printer = utils.DefaultPrinter

// GitHubIntegration 处理所有与GitHub相关的集成功能
type GitHubIntegration struct {
	config   *githubModels.GitHubConfig
	mu       sync.RWMutex
	eventBus *events.EventBus
	repos    []githubModels.GitHubRepo
	logger   *logger.GitHubLogger
}

// NewGitHubIntegration 创建一个新的GitHub集成实例
// eventBus: 用于发布集成相关的事件
func NewGitHubIntegration(eventBus *events.EventBus) *GitHubIntegration {
	return &GitHubIntegration{
		config:   &githubModels.GitHubConfig{Installations: make(map[int64]*githubModels.Installation)},
		eventBus: eventBus,
		repos:    make([]githubModels.GitHubRepo, 0),
		logger:   logger.NewGitHubLogger(),
	}
}

// HandleSetup 处理GitHub App的安装设置
// 生成GitHub App manifest并返回重定向URL
// 返回值:
//   - string: 重定向URL
//   - error: 处理过程中的错误
func (g *GitHubIntegration) HandleSetup(c *gin.Context) (string, error) {
	g.logger.LogInfo("开始处理 GitHub App 安装设置")
	return config.GenerateManifest(c)
}

// HandleCallback 处理GitHub App安装后的回调
// 处理GitHub的OAuth回调，保存应用凭据
// 返回值:
//   - string: 安装URL
//   - error: 处理过程中的错误
func (g *GitHubIntegration) HandleCallback(c *gin.Context) (string, error) {
	g.logger.LogInfo("开始处理 GitHub App 安装回调")
	result, err := config.ProcessCallback(c)
	if err != nil {
		return "", err
	}

	g.mu.Lock()
	g.config.GitHubAppID = result.ID
	g.config.GitHubAppPrivateKey = result.PEM
	g.mu.Unlock()

	return result.GetInstallURL(), nil
}

// HandleWebhook 处理接收到的GitHub webhook事件
// 验证webhook签名，保存事件数据，并触发相应的处理逻辑
// 参数:
//   - c: Gin上下文
//
// 返回值:
//   - error: 处理过程中的错误
func (g *GitHubIntegration) HandleWebhook(c *gin.Context) error {
	return webhook.ProcessWebhookEvent(c, g.config, g.eventBus)
}

// GetStoredWebhooks 获取所有存储的webhook事件数据
// 返回值:
//   - []models.WebhookPayload: webhook事件列表
//   - error: 获取过程中的错误
func (g *GitHubIntegration) GetStoredWebhooks() ([]githubModels.WebhookPayload, error) {
	return storage.GetWebhooks(storage.WebhookDir)
}

// GetConfig 返回当前的GitHub配置信息
// 返回值:
//   - *models.GitHubConfig: GitHub配置信息
func (g *GitHubIntegration) GetConfig() *githubModels.GitHubConfig {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.config
}

// ListAuthorizedRepos 获取已授权的仓库列表
// 返回值:
//   - []GitHubRepo: 已授权的仓库列表
//   - error: 获取过程中的错误
func (g *GitHubIntegration) ListAuthorizedRepos(ctx context.Context) ([]githubModels.GitHubRepo, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.repos, nil
}

// GetLogs 获取GitHub集成的日志目录内容
func (g *GitHubIntegration) GetLogs() ([]models.FileInfo, error) {
	return g.logger.GetLogFiles()
}
