package controllers

import (
	"fmt"
	"net/http"
	"servon/core/internal/libs/github"
	"servon/core/internal/libs/github/models"

	"github.com/gin-gonic/gin"
)

type GitHubController struct {
	*github.GitHubIntegration
}

func NewGitHubController(integrations *github.GitHubIntegration) *GitHubController {
	return &GitHubController{GitHubIntegration: integrations}
}

// HandleGitHubSetup handles GitHub App Manifest flow setup request
func (h *GitHubController) HandleGitHubSetup(c *gin.Context) {
	var req struct {
		BaseURL     string `json:"base_url" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %v", err),
		})
		return
	}

	manifest, err := h.GitHubIntegration.HandleSetup(req.Name, req.Description, req.BaseURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("failed to setup GitHub App: %v", err),
		})
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, manifest)
}

// HandleGitHubCallback handles callback after GitHub App creation
func (h *GitHubController) HandleGitHubCallback(c *gin.Context) {
	printer.PrintInfof("HandleGitHubCallback")
	redirectURL, err := h.GitHubIntegration.HandleCallback(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// HandleGitHubWebhook 处理来自 GitHub 的 webhook 请求
func (h *GitHubController) HandleGitHubWebhook(c *gin.Context) {
	printer.PrintInfof("HandleGitHubWebhook")
	if err := h.GitHubIntegration.HandleWebhook(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// HandleGetWebhooks 获取存储的 webhook 数据
func (h *GitHubController) HandleGetWebhooks(c *gin.Context) {
	printer.PrintInfof("HandleGetWebhooks")
	webhooks, err := h.GitHubIntegration.GetStoredWebhooks()
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
