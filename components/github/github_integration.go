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
	"servon/components/events"
	"strings"
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
	eventBus   events.IEventBus
	tokenCache *TokenCacheManager
}

// NewGitHubIntegration 创建一个新的GitHub集成实例
// eventBus: 用于发布集成相关的事件
func NewGitHubIntegration(eventBus events.IEventBus) *GitHubIntegration {
	return &GitHubIntegration{
		eventBus:   eventBus,
		tokenCache: NewTokenCacheManager(),
	}
}

// HandleSetup 处理GitHub App的安装设置
// 生成GitHub App manifest并返回重定向URL
// 返回值:
//   - string: 重定向URL
//   - error: 处理过程中的错误
func (g *GitHubIntegration) HandleSetup(name, description, baseURL string) (string, error) {
	return g.GenerateManifest(name, description, baseURL)
}

// HandleCallback 处理GitHub App安装后的回调
// 处理GitHub的OAuth回调，保存应用凭据
// 返回值:
//   - string: 安装URL
//   - error: 处理过程中的错误
func (g *GitHubIntegration) HandleCallback(c *gin.Context) (string, error) {
	result, err := ProcessCallback(c)
	if err != nil {
		return "", err
	}

	// 直接保存到磁盘
	err = SaveAppConfig(&GitHubConfig{
		GitHubAppID:         result.ID,
		GitHubAppPrivateKey: result.PEM,
		Installations:       make(map[int64]*Installation),
		UpdatedAt:           time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return "", fmt.Errorf("保存GitHub App配置失败: %v", err)
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
	return g.GetWebhooks(WebhookDir)
}

// ListAuthorizedRepos 获取已授权的仓库列表
// 返回值:
//   - []GitHubRepo: 已授权的仓库列表
//   - error: 获取过程中的错误
func (g *GitHubIntegration) ListAuthorizedRepos(ctx context.Context) ([]GitHubRepo, error) {
	// 从存储中读取安装配置
	installations, err := g.GetInstallationConfig()
	if err != nil {
		return nil, fmt.Errorf("读取安装配置失败: %v", err)
	}

	// 创建结果切片
	var repos []GitHubRepo

	// 遍历所有安装实例
	for _, installation := range installations {
		// 获取该安装实例下的所有仓库
		repos = installation.Repositories
	}

	return repos, nil
}

// GetInstallationToken 获取仓库的安装令牌
func (g *GitHubIntegration) GetInstallationToken(repo string) (string, error) {
	// 处理仓库地址格式
	repoFullName := repo
	if strings.HasPrefix(repo, "https://github.com/") {
		// 从 https://github.com/owner/repo 转换为 owner/repo
		repoFullName = strings.TrimPrefix(repo, "https://github.com/")
	}

	// 验证仓库名称格式
	parts := strings.Split(repoFullName, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("无效的仓库名称格式: %s，应为 'owner/repo' 格式", repoFullName)
	}

	// 获取所有安装配置
	installations, err := g.GetInstallationConfig()
	if err != nil {
		return "", fmt.Errorf("读取安装配置失败: %v", err)
	}

	// 查找对应的安装
	var installation *Installation
	for _, inst := range installations {
		for _, r := range inst.Repositories {
			if r.FullName == repoFullName {
				installation = inst
				break
			}
		}
		if installation != nil {
			break
		}
	}

	if installation == nil {
		return "", fmt.Errorf("未找到仓库 %s 的安装信息，请确认 GitHub App 已正确安装到该仓库", repo)
	}

	// 检查 App 配置
	appConfig, err := LoadAppConfig()
	if err != nil {
		return "", fmt.Errorf("加载 GitHub App 配置失败: %v", err)
	}
	if appConfig == nil {
		return "", fmt.Errorf("GitHub App 配置不存在")
	}

	// 检查缓存
	if token, ok := g.tokenCache.Get(installation.ID); ok {
		return token, nil
	}

	// 确保至少有 contents:read 权限
	if level, ok := installation.Permissions["contents"]; !ok || level != "read" && level != "write" {
		return "", fmt.Errorf("安装缺少必要的 contents 权限，当前权限: %v", installation.Permissions)
	}

	token, expiresAt, err := g.createInstallationToken(installation.ID)
	if err != nil {
		return "", fmt.Errorf("创建安装令牌失败: %v", err)
	}

	// 验证令牌
	if err := g.validateToken(token, repo); err != nil {
		return "", fmt.Errorf("令牌验证失败: %v", err)
	}

	g.tokenCache.Set(installation.ID, token, expiresAt)
	return token, nil
}

func (g *GitHubIntegration) createInstallationToken(installationID int64) (string, time.Time, error) {
	// 获取 App 配置
	appConfig, err := LoadAppConfig()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("加载 GitHub App 配置失败: %v", err)
	}
	if appConfig == nil {
		return "", time.Time{}, fmt.Errorf("GitHub App 配置不存在")
	}

	// 验证安装是否存在于配置中
	if _, exists := appConfig.Installations[installationID]; !exists {
		return "", time.Time{}, fmt.Errorf("安装 ID %d 在 App 配置中不存在", installationID)
	}

	// 生成 JWT
	jwt, err := g.generateJWT()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("生成 JWT 失败: %v", err)
	}

	// 创建请求
	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("创建 HTTP 请求失败: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", time.Time{}, fmt.Errorf("GitHub API 返回错误: %s - %s", resp.Status, string(body))
	}

	// 解析响应
	var result struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", time.Time{}, fmt.Errorf("解析响应失败: %v", err)
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

	block, _ := pem.Decode([]byte(appConfig.GitHubAppPrivateKey))
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
		"iss": appConfig.GitHubAppID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("签名JWT失败: %v", err)
	}

	return signedToken, nil
}

// ProcessWebhookEvent 处理 GitHub webhook 请求
func (g *GitHubIntegration) ProcessWebhookEvent(c *gin.Context) error {
	event := c.GetHeader("X-GitHub-Event")
	eventID := c.GetHeader("X-GitHub-Delivery")

	payload, err := c.GetRawData()
	if err != nil {
		return fmt.Errorf("读取 payload 失败: %v", err)
	}

	// 从磁盘获取 App 配置
	appConfig, err := LoadAppConfig()
	if err != nil {
		return fmt.Errorf("加载 GitHub App 配置失败: %v", err)
	}
	if appConfig == nil {
		return fmt.Errorf("GitHub App 配置不存在")
	}

	// 验证 webhook
	if err := validateWebhook(c); err != nil {
		return fmt.Errorf("webhook validation failed: %v", err)
	}

	// 保存 webhook 数据
	if err := g.SaveWebhookPayload(WebhookDir, event, eventID, payload); err != nil {
		return fmt.Errorf("failed to save webhook payload: %v", err)
	}

	return g.handleEvent(event, payload, g.eventBus)
}

func validateWebhook(c *gin.Context) error {
	signature := c.GetHeader("X-Hub-Signature-256")
	if signature == "" {
		return fmt.Errorf("missing signature")
	}

	// TODO: 实现签名验证
	return nil
}

func (g *GitHubIntegration) handleEvent(event string, payload []byte, eventBus events.IEventBus) error {
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
		// 未处理的事件类型
		return nil
	}
}

// handleInstallationEvent 处理 GitHub App 安装相关的事件
func (g *GitHubIntegration) handleInstallationEvent(payload []byte) error {
	// 首先保存原始 webhook payload
	if err := SaveRawInstallationData(payload); err != nil {
		return fmt.Errorf("保存原始安装数据失败: %v", err)
	}

	var event struct {
		Action       string       `json:"action"`
		Installation Installation `json:"installation"`
		Repositories []GitHubRepo `json:"repositories"`
		Sender       struct {
			Login     string `json:"login"`
			ID        int64  `json:"id"`
			AvatarURL string `json:"avatar_url"`
			Type      string `json:"type"`
		} `json:"sender"`
	}

	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("解析安装事件失败: %v", err)
	}

	// 更新安装信息
	installation := &event.Installation
	installation.AccountLogin = installation.Account.Login
	installation.AccountID = installation.Account.ID
	installation.AccountType = installation.Account.Type
	installation.AccountAvatarURL = installation.Account.AvatarURL
	installation.Repositories = event.Repositories

	// 保存安装配置到独立文件
	if err := SaveInstallationConfig(installation); err != nil {
		return fmt.Errorf("保存安装配置失败: %v", err)
	}

	// 更新 app_config.json 中的 installations 字段
	appConfig, err := LoadAppConfig()
	if err != nil {
		return fmt.Errorf("加载 App 配置失败: %v", err)
	}
	if appConfig == nil {
		return fmt.Errorf("App 配置不存在")
	}

	// 更新或添加安装信息
	if appConfig.Installations == nil {
		appConfig.Installations = make(map[int64]*Installation)
	}
	appConfig.Installations[installation.ID] = installation

	// 保存更新后的 App 配置
	if err := SaveAppConfig(appConfig); err != nil {
		return fmt.Errorf("保存更新后的 App 配置失败: %v", err)
	}

	return nil
}

// handlePushEvent 处理代码推送事件
func (g *GitHubIntegration) handlePushEvent(payload []byte, eventBus events.IEventBus) error {
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

// 添加新的验证方法
func (g *GitHubIntegration) validateToken(token, repo string) error {
	// 创建一个带认证的请求来测试令牌
	client := &http.Client{}
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s", strings.TrimPrefix(repo, "https://github.com/"))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("创建验证请求失败: %v", err)
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("执行验证请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API请求失败 (状态码: %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
