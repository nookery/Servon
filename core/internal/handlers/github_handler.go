package handlers

import (
	"net/http"

	"servon/core/internal/libs/github/config"
	"servon/core/internal/libs/github/storage"
	"servon/core/internal/libs/github/webhook"

	"github.com/gin-gonic/gin"
)

var (
	githubConfig   = config.NewGitHubConfig()
	webhookHandler = webhook.NewHandler(githubConfig)
	storageManager = storage.NewManager()
)

// HandleGitHubSetup handles GitHub App Manifest flow setup request
func HandleGitHubSetup(c *gin.Context) {
	manifest, err := githubConfig.HandleSetup(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, manifest)
}

// HandleGitHubCallback handles callback after GitHub App creation
func HandleGitHubCallback(c *gin.Context) {
	redirectURL, err := githubConfig.HandleCallback(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// HandleGitHubWebhook handles webhook requests from GitHub
func HandleGitHubWebhook(c *gin.Context) {
	if err := webhookHandler.Handle(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// HandleGetWebhooks retrieves stored webhook data
func HandleGetWebhooks(c *gin.Context) {
	webhooks, err := storageManager.GetWebhooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webhooks)
}
