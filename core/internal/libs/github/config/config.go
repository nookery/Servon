// Package config 提供了 GitHub App 配置相关的功能
// 主要负责：
// 1. 生成 GitHub App Manifest
// 2. 处理 GitHub App 的回调
// 3. 管理 GitHub App 的配置信息
package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"servon/core/internal/libs/github/models"

	"github.com/gin-gonic/gin"
)

// GenerateManifest 生成 GitHub App 的 manifest
// 接收 gin.Context 作为参数，从中获取应用名称、描述和基础URL
// 返回包含 manifest 的 HTML 表单，用于自动提交到 GitHub
func GenerateManifest(c *gin.Context) (string, error) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		BaseURL     string `json:"base_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		return "", fmt.Errorf("invalid request: %v", err)
	}

	manifest := createManifest(req.Name, req.Description, req.BaseURL)
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return "", fmt.Errorf("failed to generate manifest: %v", err)
	}

	state, err := generateState()
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %v", err)
	}

	return generateHTML(state, string(manifestJSON)), nil
}

// ProcessCallback 处理 GitHub App 创建后的回调
// 从回调请求中获取 code，并使用该 code 获取 GitHub App 的详细信息
// 返回 AppCreationResult，包含应用ID、名称和私钥
func ProcessCallback(c *gin.Context) (*AppCreationResult, error) {
	code := c.Query("code")
	if code == "" {
		return nil, fmt.Errorf("missing code parameter")
	}

	return convertManifest(code)
}

// AppCreationResult 表示 GitHub App 创建的结果
type AppCreationResult struct {
	ID   int64  `json:"id"`   // GitHub App 的唯一标识符
	Name string `json:"name"` // GitHub App 的名称
	PEM  string `json:"pem"`  // GitHub App 的私钥
}

// GetInstallURL 返回 GitHub App 的安装URL
func (r *AppCreationResult) GetInstallURL() string {
	return fmt.Sprintf("https://github.com/apps/%s/installations/new", r.Name)
}

// Private helper functions

// createManifest 创建 GitHub App 的 manifest 配置
// 包含应用的基本信息、权限和事件订阅
func createManifest(name, description, baseURL string) models.GitHubManifest {
	manifest := models.GitHubManifest{
		Name:        name,
		URL:         baseURL,
		Description: description,
		Public:      true,
		HookAttributes: struct {
			URL    string `json:"url"`
			Active bool   `json:"active"`
		}{
			URL:    fmt.Sprintf("%s/web_api/github/webhook", baseURL),
			Active: true,
		},
		RedirectURL:  fmt.Sprintf("%s/web_api/github/callback", baseURL),
		CallbackURLs: []string{fmt.Sprintf("%s/web_api/github/callback", baseURL)},
		DefaultPermissions: map[string]string{
			"issues": "write",
			"checks": "write",
		},
		DefaultEvents: []string{
			"issues",
			"issue_comment",
			"check_suite",
			"check_run",
		},
	}
	return manifest
}

// generateState 生成随机的状态字符串，用于防止CSRF攻击
func generateState() (string, error) {
	state := make([]byte, 16)
	if _, err := rand.Read(state); err != nil {
		return "", err
	}
	return hex.EncodeToString(state), nil
}

// generateHTML 生成包含 manifest 的 HTML 表单
func generateHTML(state, manifestJSON string) string {
	return fmt.Sprintf(`
		<form id="github-form" action="https://github.com/settings/apps/new?state=%s" method="post">
			<input type="hidden" name="manifest" value='%s'>
		</form>
		<script>document.getElementById("github-form").submit();</script>
	`, state, manifestJSON)
}

// convertManifest 将 manifest 转换为实际的 GitHub App
// 使用 GitHub API 完成转换过程
func convertManifest(code string) (*AppCreationResult, error) {
	resp, err := http.Post(
		fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions", code),
		"application/json",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub App: %v", err)
	}
	defer resp.Body.Close()

	var result AppCreationResult

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub response: %v", err)
	}

	return &result, nil
}
