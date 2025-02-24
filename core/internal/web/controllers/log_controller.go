package controllers

import (
	"net/http"
	"strconv"

	"servon/core/internal/managers"

	"github.com/gin-gonic/gin"
)

type LogController struct {
	logManager *managers.LogManager
}

func NewLogController(logManager *managers.LogManager) *LogController {
	return &LogController{
		logManager: logManager,
	}
}

// HandleListLogFiles 获取日志文件列表
func (c *LogController) HandleListLogFiles(ctx *gin.Context) {
	subDir := ctx.DefaultQuery("dir", "")
	files, err := c.logManager.ListLogFiles(subDir)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"files": files})
}

// HandleReadLogEntries 读取日志内容
func (c *LogController) HandleReadLogEntries(ctx *gin.Context) {
	logFile := ctx.Query("file")
	if logFile == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file parameter is required"})
		return
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "100"))
	entries, err := c.logManager.ReadLogEntries(logFile, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"entries": entries})
}

// HandleSearchLogs 搜索日志
func (c *LogController) HandleSearchLogs(ctx *gin.Context) {
	subDir := ctx.DefaultQuery("dir", "")
	keyword := ctx.Query("keyword")
	if keyword == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "keyword parameter is required"})
		return
	}

	entries, err := c.logManager.SearchLogs(subDir, keyword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"entries": entries})
}

// HandleGetLogStats 获取日志统计信息
func (c *LogController) HandleGetLogStats(ctx *gin.Context) {
	subDir := ctx.DefaultQuery("dir", "")
	stats, err := c.logManager.GetLogStats(subDir)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"stats": stats})
}

// HandleCleanOldLogs 清理旧日志
func (c *LogController) HandleCleanOldLogs(ctx *gin.Context) {
	days, err := strconv.Atoi(ctx.DefaultQuery("days", "30"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid days parameter"})
		return
	}

	if err := c.logManager.CleanOldLogs(days); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Old logs cleaned successfully"})
}
