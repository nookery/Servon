package handlers

import (
	"net/http"

	"servon/core/internal/managers"

	"github.com/gin-gonic/gin"
)

var githubManager = managers.DefaultGitHubManager

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

// HandleGitHubWebhook handles webhook requests from GitHub
func HandleGitHubWebhook(c *gin.Context) {
	if err := githubManager.HandleWebhook(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// HandleGetWebhooks retrieves stored webhook data
func HandleGetWebhooks(c *gin.Context) {
	webhooks, err := githubManager.GetStoredWebhooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webhooks)
}
