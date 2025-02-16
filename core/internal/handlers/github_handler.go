package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var printer = DefaultPrinter
var config = &GitHubConfig{}

// GitHub App 集成流程:
// 1. 用户调用 /web_api/github/setup 接口，提供 App 名称和描述
// 2. 系统生成 GitHub App Manifest，重定向到 GitHub 创建页面
// 3. 用户在 GitHub 确认创建后，GitHub 回调 /web_api/github/callback
// 4. 系统保存 App 配置，并重定向用户到 App 安装页面
// 5. 用户选择要安装 App 的仓库
// 6. GitHub 发送 installation 事件到 webhook
// 7. 后续用户仓库的 push、PR 等操作都会触发相应的 webhook 事件

// 添加新的配置和存储结构
type Installation struct {
	ID           int64    `json:"id"`
	AccountID    int64    `json:"account_id"`
	AccountLogin string   `json:"account_login"`
	Repositories []string `json:"repositories"` // 仓库全名列表
}

type GitHubConfig struct {
	// GitHub integration settings
	GitHubAppID         int64  `json:"github_app_id"`
	GitHubAppPrivateKey string `json:"github_app_private_key"`
	GitHubWebhookSecret string `json:"github_webhook_secret"`

	// 新增字段用于存储安装信息
	Installations map[int64]*Installation `json:"installations"` // key: installation_id
}

func init() {
	config.Installations = make(map[int64]*Installation)
}

// GitHubManifest represents the GitHub App manifest configuration
type GitHubManifest struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	HookAttributes struct {
		URL    string `json:"url"`
		Active bool   `json:"active"`
	} `json:"hook_attributes"`
	RedirectURL        string            `json:"redirect_url"`
	CallbackURLs       []string          `json:"callback_urls"`
	Description        string            `json:"description"`
	Public             bool              `json:"public"`
	DefaultEvents      []string          `json:"default_events"`
	DefaultPermissions map[string]string `json:"default_permissions"`
}

// HandleGitHubSetup 处理GitHub App Manifest flow的设置请求
func HandleGitHubSetup(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供必要的GitHub App信息"})
		return
	}

	baseURL := "http://43.142.208.212:9754"

	// 构建manifest
	manifest := GitHubManifest{
		Name:        req.Name,
		URL:         baseURL,
		Description: req.Description,
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

	// 将manifest转换为JSON
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成manifest失败"})
		return
	}

	// 生成state参数
	state := make([]byte, 16)
	if _, err := rand.Read(state); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成state失败"})
		return
	}
	stateStr := hex.EncodeToString(state)

	// 生成HTML表单
	html := fmt.Sprintf(`
		<form id="github-form" action="https://github.com/settings/apps/new?state=%s" method="post">
			<input type="hidden" name="manifest" value='%s'>
		</form>
		<script>document.getElementById("github-form").submit();</script>
	`, stateStr, string(manifestJSON))

	// 返回HTML
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// HandleGitHubCallback 处理GitHub App创建后的回调
func HandleGitHubCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code parameter"})
		return
	}

	// 调用GitHub API完成App创建
	resp, err := http.Post(
		fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions", code),
		"application/json",
		nil,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create GitHub App: %v", err)})
		return
	}
	defer resp.Body.Close()

	var result struct {
		ID         int64  `json:"id"`
		Name       string `json:"name"`
		WebhookURL string `json:"webhook_url"`
		PEM        string `json:"pem"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse GitHub response: %v", err)})
		return
	}

	// 保存App配置
	config.GitHubAppID = result.ID
	config.GitHubAppPrivateKey = result.PEM

	// TODO: 将配置保存到持久化存储
	if err := saveConfig(); err != nil {
		printer.PrintError(err)
	}

	// 重定向到GitHub App安装页面
	installURL := fmt.Sprintf("https://github.com/apps/%s/installations/new", result.Name)
	c.Redirect(http.StatusTemporaryRedirect, installURL)
}

// HandleGitHubWebhook 处理来自GitHub的webhook请求
func HandleGitHubWebhook(c *gin.Context) {
	// 验证webhook签名
	signature := c.GetHeader("X-Hub-Signature-256")
	if signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing signature"})
		return
	}

	// 读取请求体
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read payload"})
		return
	}

	// 验证webhook secret
	// TODO: 实现签名验证

	// 解析事件类型
	event := c.GetHeader("X-GitHub-Event")
	switch event {
	case "installation", "installation_repositories":
		var installEvent struct {
			Action       string `json:"action"`
			Installation struct {
				ID      int64 `json:"id"`
				Account struct {
					ID    int64  `json:"id"`
					Login string `json:"login"`
				} `json:"account"`
			} `json:"installation"`
			Repositories []struct {
				FullName string `json:"full_name"`
			} `json:"repositories"`
		}
		if err := json.Unmarshal(payload, &installEvent); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid installation event payload"})
			return
		}

		switch installEvent.Action {
		case "created", "added":
			// 新安装或添加仓库
			installation, exists := config.Installations[installEvent.Installation.ID]
			if !exists {
				installation = &Installation{
					ID:           installEvent.Installation.ID,
					AccountID:    installEvent.Installation.Account.ID,
					AccountLogin: installEvent.Installation.Account.Login,
					Repositories: make([]string, 0),
				}
				config.Installations[installEvent.Installation.ID] = installation
			}

			// 更新仓库列表
			for _, repo := range installEvent.Repositories {
				installation.Repositories = append(installation.Repositories, repo.FullName)
			}

			printer.PrintInfo(fmt.Sprintf("Installation %d: Added repositories for %s",
				installation.ID, installation.AccountLogin))

		case "removed":
			// 移除仓库
			if installation, exists := config.Installations[installEvent.Installation.ID]; exists {
				removedRepos := make(map[string]bool)
				for _, repo := range installEvent.Repositories {
					removedRepos[repo.FullName] = true
				}

				// 过滤出未被移除的仓库
				remaining := make([]string, 0)
				for _, repo := range installation.Repositories {
					if !removedRepos[repo] {
						remaining = append(remaining, repo)
					}
				}
				installation.Repositories = remaining

				printer.PrintInfo(fmt.Sprintf("Installation %d: Removed repositories for %s",
					installation.ID, installation.AccountLogin))
			}

		case "deleted":
			// 完全卸载
			delete(config.Installations, installEvent.Installation.ID)
			printer.PrintInfo(fmt.Sprintf("Installation %d: Deleted for %s",
				installEvent.Installation.ID, installEvent.Installation.Account.Login))
		}

		// 保存更新后的配置
		if err := saveConfig(); err != nil {
			printer.PrintError(err)
		}

	case "push":
		// 处理push事件
		var pushEvent struct {
			Ref        string `json:"ref"`
			Before     string `json:"before"`
			After      string `json:"after"`
			Created    bool   `json:"created"`
			Deleted    bool   `json:"deleted"`
			Repository struct {
				FullName string `json:"full_name"`
			} `json:"repository"`
		}
		if err := json.Unmarshal(payload, &pushEvent); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid push event payload"})
			return
		}

		printer.PrintInfo(fmt.Sprintf("Received push event for %s", pushEvent.Repository.FullName))
		// TODO: 处理push事件的具体逻辑

	case "pull_request":
		// 处理pull request事件
		var prEvent struct {
			Action      string `json:"action"`
			PullRequest struct {
				Number int    `json:"number"`
				Title  string `json:"title"`
				State  string `json:"state"`
			} `json:"pull_request"`
			Repository struct {
				FullName string `json:"full_name"`
			} `json:"repository"`
		}
		if err := json.Unmarshal(payload, &prEvent); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pull request event payload"})
			return
		}

		printer.PrintInfo(fmt.Sprintf("Received PR event for %s: %s", prEvent.Repository.FullName, prEvent.Action))
		// TODO: 处理PR事件的具体逻辑

	default:
		printer.PrintInfo(fmt.Sprintf("Received unhandled event type: %s", event))
	}

	c.Status(http.StatusOK)
}

// saveConfig 保存配置到持久化存储
func saveConfig() error {
	// TODO: 实现配置的持久化存储
	// 可以保存到文件或数据库
	return nil
}
