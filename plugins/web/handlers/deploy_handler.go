package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"servon/core"

	"github.com/gin-gonic/gin"
)

const deployLogsPath = "/data/deploy"

func init() {
	// 确保部署日志目录存在
	if err := os.MkdirAll(deployLogsPath, 0755); err != nil {
		panic("无法创建部署日志目录: " + err.Error())
	}
}

// HandleGetDeployLog 获取单个部署日志
func (h *WebHandler) HandleGetDeployLog(c *gin.Context) {
	logID := c.Query("id")
	if logID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需要提供日志ID"})
		return
	}

	filePath := filepath.Join(deployLogsPath, logID+".json")
	content, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "读取日志失败: " + err.Error()})
		return
	}

	var log core.DeployLog
	if err := json.Unmarshal(content, &log); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析日志失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, log)
}

// HandleListDeployLogs 获取部署日志列表
func (h *WebHandler) HandleListDeployLogs(c *gin.Context) {
	files, err := os.ReadDir(deployLogsPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取日志目录失败: " + err.Error()})
		return
	}

	var logs []core.DeployLog
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		content, err := os.ReadFile(filepath.Join(deployLogsPath, file.Name()))
		if err != nil {
			continue
		}

		var log core.DeployLog
		if err := json.Unmarshal(content, &log); err != nil {
			continue
		}
		logs = append(logs, log)
	}

	c.JSON(http.StatusOK, logs)
}

// HandleDeleteDeployLog 删除部署日志
func (h *WebHandler) HandleDeleteDeployLog(c *gin.Context) {
	logID := c.Query("id")
	if logID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需要提供日志ID"})
		return
	}

	filePath := filepath.Join(deployLogsPath, logID+".json")
	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除日志失败: " + err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
