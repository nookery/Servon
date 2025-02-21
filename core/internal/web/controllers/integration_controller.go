package controllers

import (
	"fmt"
	"net/http"
	"servon/core/internal/managers"
	"servon/core/internal/managers/github"

	"github.com/gin-gonic/gin"
)

type IntegrationController struct {
	*managers.FullManager
}

func NewIntegrationController(fullIntegration *managers.FullManager) *IntegrationController {
	return &IntegrationController{FullManager: fullIntegration}
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
