package controllers

import (
	"fmt"
	"net/http"
	"servon/core/internal/integrations/github"

	"github.com/gin-gonic/gin"
)

type IntegrationController struct {
	*github.GitHubIntegration
}

func NewIntegrationController(fullIntegration *github.GitHubIntegration) *IntegrationController {
	return &IntegrationController{GitHubIntegration: fullIntegration}
}

// HandleListGitHubRepos 获取GitHub授权仓库列表
func (h *IntegrationController) HandleListGitHubRepos(c *gin.Context) {
	repos, err := h.ListAuthorizedRepos(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取GitHub仓库列表失败: %v", err),
		})
		return
	}

	if repos == nil {
		repos = make([]github.GitHubRepo, 0)
	}

	c.JSON(http.StatusOK, repos)
}

// HandleGetGitHubLogs 获取GitHub集成日志
func (h *IntegrationController) HandleGetGitHubLogs(c *gin.Context) {
	logs, err := h.GetLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取GitHub日志失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, logs)
}
