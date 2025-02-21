// Package github 提供了与GitHub App集成的核心功能实现
//
// GitHub App 集成流程:
//
//  1. 安装流程
//     - 调用 HandleSetup 生成GitHub App manifest
//     - GitHub重定向用户到manifest配置页面
//     - 用户确认后GitHub回调HandleCallback
//     - 保存GitHub App凭据(ID和私钥)
//     - 重定向用户到GitHub App安装页面
//
//  2. Webhook配置
//     - GitHub App安装后自动配置webhook
//     - HandleWebhook接收和处理webhook事件
//     - 事件通过EventBus分发到系统其他部分
//
//  3. 认证机制
//     - 使用GitHub App私钥生成JWT
//     - 使用JWT获取Installation Token
//     - Token缓存管理避免频繁请求
//     - 自动处理Token过期和更新
//
// 主要功能:
//   - GitHub App的安装和配置
//   - Webhook事件的接收和处理
//   - 安装令牌的管理和缓存
//   - 仓库访问权限的控制
//
// 使用示例:
//
//  1. 安装GitHub App:
//     ```go
//     integration := NewGitHubIntegration(eventBus)
//     redirectURL, err := integration.HandleSetup("My App", "Description", "https://example.com")
//     // 重定向用户到GitHub
//     ```
//
//  2. 获取仓库访问令牌:
//     ```go
//     token, err := integration.GetInstallationToken("owner/repo")
//     if err != nil {
//     log.Printf("获取令牌失败: %v", err)
//     return
//     }
//     // 使用token访问仓库
//     ```
//
// 安全注意事项:
//  1. GitHub App私钥必须安全存储
//  2. Webhook密钥需要定期轮换
//  3. 安装令牌应该及时清理
//  4. 权限应该遵循最小权限原则
//
// 错误处理:
//   - 所有公开方法都返回详细的错误信息
//   - 包含重试机制处理临时性故障
//   - 记录关键操作的日志
//
// 配置要求:
//
//  1. GitHub App权限:
//     - Repository contents: Read (用于访问代码)
//     - Metadata: Read (用于访问仓库信息)
//     - Webhooks: Read & Write (用于配置webhook)
//
//  2. 系统配置:
//     - 需要配置webhook接收URL
//     - 需要配置GitHub App的ID和私钥
//     - 需要配置webhook密钥
package github

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"servon/core/internal/events"
	"servon/core/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// GitHubIntegration 处理所有与GitHub相关的集成功能
// 主要职责:
//  1. 管理GitHub App的安装和配置
//  2. 处理webhook事件
//  3. 管理访问令牌
//  4. 提供仓库访问接口
type GitHubIntegration struct {
	eventBus   *events.EventBus
	logger     *GitHubLogger
	tokenCache *TokenCacheManager
}

// NewGitHubIntegration 创建一个新的GitHub集成实例
// eventBus: 用于发布集成相关的事件
func NewGitHubIntegration(eventBus *events.EventBus) *GitHubIntegration {
	return &GitHubIntegration{
		eventBus:   eventBus,
		logger:     DefaultGitHubLogger,
		tokenCache: NewTokenCacheManager(),
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

	// 直接保存到磁盘
	err = SaveAppConfig(&AppConfig{
		AppID:      result.ID,
		PrivateKey: result.PEM,
	})
	if err != nil {
		g.logger.LogErrorf("保存GitHub App配置失败: %v", err)
		return "", err // 这里应该返回错误，因为没有内存配置可用
	}

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
	return g.ProcessWebhookEvent(c)
}

// GetStoredWebhooks 获取所有存储的webhook事件数据
// 返回值:
//   - []WebhookPayload: webhook事件列表
//   - error: 获取过程中的错误
func (g *GitHubIntegration) GetStoredWebhooks() ([]WebhookPayload, error) {
	return GetWebhooks(WebhookDir)
}

// ListAuthorizedRepos 获取已授权的仓库列表
// 返回值:
//   - []GitHubRepo: 已授权的仓库列表
//   - error: 获取过程中的错误
func (g *GitHubIntegration) ListAuthorizedRepos(ctx context.Context) ([]GitHubRepo, error) {
	// 从存储中读取安装配置
	installations, err := GetInstallationConfig()
	if err != nil {
		g.logger.LogErrorf("读取安装配置失败: %v", err)
		return nil, err
	}

	// 创建结果切片
	var repos []GitHubRepo

	// 遍历所有安装实例
	for _, installation := range installations {
		g.logger.LogInfof("正在处理安装实例: %v", installation)
		// 获取该安装实例下的所有仓库
		for _, repoName := range installation.Repositories {
			repo := GitHubRepo{
				Name:     repoName.Name,
				FullName: fmt.Sprintf("%s/%s", installation.AccountLogin, repoName.Name),
				Private:  repoName.Private,
				HTMLURL:  fmt.Sprintf("https://github.com/%s/%s", installation.AccountLogin, repoName.FullName),
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

// GetInstallationToken 获取安装令牌，支持缓存
func (g *GitHubIntegration) GetInstallationToken(repo string) (string, error) {
	// 从存储中读取安装信息
	installations, err := GetInstallationConfig()
	if err != nil {
		return "", fmt.Errorf("读取安装配置失败: %v", err)
	}

	// 查找对应的安装实例
	var installation *Installation
	for _, inst := range installations {
		for _, r := range inst.Repositories {
			if r.FullName == repo {
				installation = inst
				break
			}
		}
		if installation != nil {
			break
		}
	}

	if installation == nil {
		return "", fmt.Errorf("未找到仓库 %s 的安装信息", repo)
	}

	// 检查缓存
	if token, ok := g.tokenCache.Get(installation.ID); ok {
		return token, nil
	}

	// 缓存不存在或已过期，创建新令牌
	token, expiresAt, err := g.createInstallationToken(installation.ID)
	if err != nil {
		return "", fmt.Errorf("创建安装令牌失败: %v", err)
	}

	// 更新缓存
	g.tokenCache.Set(installation.ID, token, expiresAt)

	return token, nil
}

// createInstallationToken 创建新的安装令牌
func (g *GitHubIntegration) createInstallationToken(installationID int64) (string, time.Time, error) {
	// 生成 JWT
	jwt, err := g.generateJWT()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("生成 JWT 失败: %v", err)
	}

	// 准备请求
	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", time.Time{}, err
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", time.Time{}, fmt.Errorf("GitHub API 错误: %s - %s", resp.Status, string(body))
	}

	// 解析响应
	var result struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", time.Time{}, err
	}

	return result.Token, result.ExpiresAt, nil
}

// generateJWT 生成用于GitHub API认证的JWT
func (g *GitHubIntegration) generateJWT() (string, error) {
	appConfig, err := LoadAppConfig()
	if err != nil {
		return "", fmt.Errorf("加载 GitHub App 配置失败: %v", err)
	}
	if appConfig == nil {
		return "", fmt.Errorf("GitHub App 配置不存在")
	}

	block, _ := pem.Decode([]byte(appConfig.PrivateKey))
	if block == nil {
		return "", fmt.Errorf("解析私钥失败")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("解析RSA私钥失败: %v", err)
	}

	// 创建JWT
	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(10 * time.Minute).Unix(),
		"iss": appConfig.AppID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("签名JWT失败: %v", err)
	}

	return signedToken, nil
}
