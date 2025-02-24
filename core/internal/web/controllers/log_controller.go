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

	// 确保返回空数组而不是 null
	if files == nil {
		files = make([]string, 0)
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

	// 确保返回空数组而不是 null
	if entries == nil {
		entries = make([]managers.LogEntry, 0)
	}

	// 反转日志条目顺序，让新的日志在前面
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
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

	// 确保返回空数组而不是 null
	if entries == nil {
		entries = make([]managers.LogEntry, 0)
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

// HandleDeleteLogFile 处理删除日志文件的请求
func (c *LogController) HandleDeleteLogFile(ctx *gin.Context) {
	// 定义请求体结构
	var req struct {
		Params struct {
			Path string `json:"path"`
		} `json:"params"`
	}

	// 解析JSON请求体
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}

	if req.Params.Path == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "需要提供文件路径"})
		return
	}

	if err := c.logManager.DeleteLogFile(req.Params.Path); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "日志文件删除成功"})
}

// HandleClearLogFile 清空指定日志文件
func (c *LogController) HandleClearLogFile(ctx *gin.Context) {
	var req struct {
		Params struct {
			File string `json:"file"`
		} `json:"params"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}

	if req.Params.File == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "需要提供文件名"})
		return
	}

	err := c.logManager.ClearLogFile(req.Params.File)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "日志文件已清空"})
}
