package handlers

import (
	"net/http"

	"servon/core/internal/libs/github/models"
	"servon/core/internal/managers"
	"servon/core/internal/utils"

	"github.com/gin-gonic/gin"
)

var githubManager = managers.DefaultGitHubManager
var printer = utils.DefaultPrinter

// HandleGitHubSetup handles GitHub App Manifest flow setup request
func HandleGitHubSetup(c *gin.Context) {
	manifest, err := githubManager.HandleSetup(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, manifest)
}

// HandleGitHubCallback handles callback after GitHub App creation
func HandleGitHubCallback(c *gin.Context) {
	redirectURL, err := githubManager.HandleCallback(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// HandleGitHubWebhook 处理来自 GitHub 的 webhook 请求
func HandleGitHubWebhook(c *gin.Context) {
	if err := githubManager.HandleWebhook(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// HandleGetWebhooks 获取存储的 webhook 数据
func HandleGetWebhooks(c *gin.Context) {
	webhooks, err := githubManager.GetStoredWebhooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 确保返回空数组而不是 null
	if webhooks == nil {
		webhooks = make([]models.WebhookPayload, 0)
	}

	printer.PrintInfof("HandleGetWebhooks: %d", len(webhooks))
	c.JSON(http.StatusOK, webhooks)
}
