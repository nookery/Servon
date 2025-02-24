// Package config 提供了 GitHub App 配置相关的功能
// 主要负责：
// 1. 生成 GitHub App Manifest
// 2. 处理 GitHub App 的回调
// 3. 管理 GitHub App 的配置信息
package github

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GenerateManifest 生成 GitHub App 的 manifest
// 接收 gin.Context 作为参数，从中获取应用名称、描述和基础URL
// 返回包含 manifest 的 HTML 表单，用于自动提交到 GitHub
func (g *GitHubIntegration) GenerateManifest(name, description, baseURL string) (string, error) {
	g.logger.Infof("GenerateManifest")

	manifest := createManifest(name, description, baseURL)
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return "", g.logger.LogAndReturnErrorf("failed to generate manifest: %v", err)
	}

	state, err := generateState()
	if err != nil {
		return "", g.logger.LogAndReturnErrorf("failed to generate state: %v", err)
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

// GetInstallURL 返回 GitHub App 的安装URL
func (r *AppCreationResult) GetInstallURL() string {
	return fmt.Sprintf("https://github.com/apps/%s/installations/new", r.Name)
}

// createManifest 创建 GitHub App 的 manifest 配置
// 包含应用的基本信息、权限和事件订阅
func createManifest(name, description, baseURL string) GitHubManifest {
	manifest := GitHubManifest{
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
			"contents": "read",  // 仓库内容读取权限
			"issues":   "write", // issues 写入权限
			"checks":   "write", // checks 写入权限
			"metadata": "read",  // 元数据读取权限
		},
		DefaultEvents: []string{
			"issues",
			"issue_comment",
			"check_suite",
			"check_run",
			"push",
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
