package handlers

import (
	"net/http"
	"servon/core"

	"github.com/gin-gonic/gin"
)

// HandleGitHubSetup handles GitHub App Manifest flow setup request
func (h *WebHandler) HandleGitHubSetup(c *gin.Context) {
	manifest, err := h.App.GitHubIntegration.HandleSetup(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, manifest)
}

// HandleGitHubCallback handles callback after GitHub App creation
func (h *WebHandler) HandleGitHubCallback(c *gin.Context) {
	h.App.Printer.PrintInfof("HandleGitHubCallback")
	redirectURL, err := h.App.GitHubIntegration.HandleCallback(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// HandleGitHubWebhook 处理来自 GitHub 的 webhook 请求
func (h *WebHandler) HandleGitHubWebhook(c *gin.Context) {
	h.App.Printer.PrintInfof("HandleGitHubWebhook")
	if err := h.App.GitHubIntegration.HandleWebhook(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// HandleGetWebhooks 获取存储的 webhook 数据
func (h *WebHandler) HandleGetWebhooks(c *gin.Context) {
	h.App.Printer.PrintInfof("HandleGetWebhooks")
	webhooks, err := h.App.GitHubIntegration.GetStoredWebhooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 确保返回空数组而不是 null
	if webhooks == nil {
		webhooks = make([]core.WebhookPayload, 0)
	}

	h.App.Printer.PrintInfof("HandleGetWebhooks: %d", len(webhooks))
	c.JSON(http.StatusOK, webhooks)
}
