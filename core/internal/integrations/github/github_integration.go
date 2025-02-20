// Package github 提供了各种第三方服务的集成实现
package github

import (
	"context"
	"fmt"
	"servon/core/internal/events"
	"servon/core/internal/models"
	"sync"

	"github.com/gin-gonic/gin"
)

// GitHubIntegration 处理所有与GitHub相关的集成功能
type GitHubIntegration struct {
	config   *GitHubConfig
	mu       sync.RWMutex
	eventBus *events.EventBus
	repos    []GitHubRepo
	logger   *GitHubLogger
}

// NewGitHubIntegration 创建一个新的GitHub集成实例
// eventBus: 用于发布集成相关的事件
func NewGitHubIntegration(eventBus *events.EventBus) *GitHubIntegration {
	return &GitHubIntegration{
		config: &GitHubConfig{
			Installations: make(map[int64]*Installation),
		},
		eventBus: eventBus,
		repos:    make([]GitHubRepo, 0),
		logger:   DefaultGitHubLogger,
	}
}

// HandleSetup 处理GitHub App的安装设置
// 生成GitHub App manifest并返回重定向URL
// 返回值:
//   - string: 重定向URL
//   - error: 处理过程中的错误
func (g *GitHubIntegration) HandleSetup(name, description, baseURL string) (string, error) {
	g.logger.LogInfo("开始处理 GitHub App 安装设置")
	return GenerateManifest(name, description, baseURL)
}

// HandleCallback 处理GitHub App安装后的回调
// 处理GitHub的OAuth回调，保存应用凭据
// 返回值:
//   - string: 安装URL
//   - error: 处理过程中的错误
func (g *GitHubIntegration) HandleCallback(c *gin.Context) (string, error) {
	g.logger.LogInfo("开始处理 GitHub App 安装回调")
	result, err := ProcessCallback(c)
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
	return g.ProcessWebhookEvent(c, g.config, g.eventBus)
}

// GetStoredWebhooks 获取所有存储的webhook事件数据
// 返回值:
//   - []WebhookPayload: webhook事件列表
//   - error: 获取过程中的错误
func (g *GitHubIntegration) GetStoredWebhooks() ([]WebhookPayload, error) {
	return GetWebhooks(WebhookDir)
}

// GetConfig 返回当前的GitHub配置信息
// 返回值:
//   - *GitHubConfig: GitHub配置信息
func (g *GitHubIntegration) GetConfig() *GitHubConfig {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.config
}

// ListAuthorizedRepos 获取已授权的仓库列表
// 返回值:
//   - []GitHubRepo: 已授权的仓库列表
//   - error: 获取过程中的错误
func (g *GitHubIntegration) ListAuthorizedRepos(ctx context.Context) ([]GitHubRepo, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// 创建结果切片
	var repos []GitHubRepo

	// 遍历所有安装实例
	for _, installation := range g.config.Installations {
		g.logger.LogInfof("installation: %v", installation)
		// 获取该安装实例下的所有仓库
		for _, repoName := range installation.Repositories {
			repo := GitHubRepo{
				Name:     repoName,
				FullName: fmt.Sprintf("%s/%s", installation.AccountLogin, repoName),
				Private:  false, // 这里可以通过 GitHub API 获取详细信息
				HTMLURL:  fmt.Sprintf("https://github.com/%s/%s", installation.AccountLogin, repoName),
			}
			repos = append(repos, repo)
		}
	}

	return repos, nil
}

// GetLogs 获取GitHub集成的日志目录内容
func (g *GitHubIntegration) GetLogs() ([]models.FileInfo, error) {
	return g.logger.GetLogFiles()
}
