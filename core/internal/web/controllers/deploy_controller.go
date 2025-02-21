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

// HandleGetDeployLog 获取单个部署日志
func (h *DeployController) HandleGetDeployLog(c *gin.Context) {
	logID := c.Query("id")
	if logID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需要提供日志ID"})
		return
	}

	log, err := h.GetDeployLog(logID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "读取日志失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, log)
}

// HandleListDeployLogs 获取部署日志列表
func (h *DeployController) HandleListDeployLogs(c *gin.Context) {
	logs, err := h.ListDeployLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取日志列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// HandleDeleteDeployLog 删除部署日志
func (h *DeployController) HandleDeleteDeployLog(c *gin.Context) {
	logID := c.Query("id")
	if logID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需要提供日志ID"})
		return
	}

	if err := h.DeleteDeployLog(logID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除日志失败: " + err.Error()})
		return
	}

	c.Status(http.StatusOK)
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
