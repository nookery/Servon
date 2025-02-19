package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleListGitHubRepos 获取GitHub授权仓库列表
func (h *WebHandler) HandleListGitHubRepos(c *gin.Context) {
	repos, err := h.ListAuthorizedRepos(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取GitHub仓库列表失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, repos)
}

// HandleGetGitHubLogs 获取GitHub集成日志
func (h *WebHandler) HandleGetGitHubLogs(c *gin.Context) {
	logs, err := h.GetLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取GitHub日志失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, logs)
}
