package controllers

import (
	"net/http"
	"servon/core/internal/managers"

	"github.com/gin-gonic/gin"
)

type DeployController struct {
	*managers.FullManager
}

func NewDeployController(manager *managers.FullManager) *DeployController {
	return &DeployController{FullManager: manager}
}

// DeployRepository 部署仓库
func (h *DeployController) DeployRepository(c *gin.Context) {
	repoID := c.Query("id")
	if repoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需要提供仓库ID"})
		return
	}

	h.DeployProject(repoID)

	c.JSON(http.StatusOK, gin.H{"message": "部署仓库成功"})
}
