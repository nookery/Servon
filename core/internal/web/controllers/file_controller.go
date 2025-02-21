package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"servon/core/internal/managers"
	"servon/core/internal/utils"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	*managers.FullManager
}

func NewFileController(manager *managers.FullManager) *FileController {
	return &FileController{FullManager: manager}
}

// HandleDeleteFile 处理删除文件的请求
func (h *FileController) HandleDeleteFile(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需要提供文件路径"})
		return
	}

	err := os.Remove(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文件失败: " + err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// HandleCreateFile 处理创建新文件的请求
func (h *FileController) HandleCreateFile(c *gin.Context) {
	var req struct {
		Path string `json:"path"`
		Type string `json:"type"` // "file" 或 "directory"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}

	// 检查路径是否已存在
	if _, err := os.Stat(req.Path); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件或目录已存在: " + req.Path})
		return
	}

	// 检查父目录是否存在且可写
	parentDir := filepath.Dir(req.Path)
	if _, err := os.Stat(parentDir); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "父目录不存在: " + parentDir})
		return
	}

	if req.Type == "directory" {
		err := os.MkdirAll(req.Path, 0755)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("创建目录失败: %v (路径: %s)", err, req.Path),
			})
			return
		}
	} else {
		// 创建空文件
		f, err := os.Create(req.Path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("创建文件失败: %v (路径: %s)", err, req.Path),
			})
			return
		}
		f.Close()
	}

	c.Status(http.StatusOK)
}

// HandleFileContent 处理获取文件内容的请求
func (h *FileController) HandleFileContent(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需要提供文件路径"})
		return
	}

	content, err := os.ReadFile(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文件失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": string(content)})
}

// HandleSaveFile 处理保存文件内容的请求
func (h *FileController) HandleSaveFile(c *gin.Context) {
	var req struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	err := os.WriteFile(req.Path, []byte(req.Content), 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败: " + err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// HandleFileDownload handles file download requests
func (h *FileController) HandleFileDownload(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path parameter is required"})
		return
	}

	// Verify the file exists and is not a directory
	fileInfo, err := os.Stat(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	if fileInfo.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot download directories"})
		return
	}

	// Serve the file
	c.File(path)
}

// HandleFileList 处理获取文件列表的请求
func (h *FileController) HandleFileList(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		path = "/"
	}

	// 获取排序参数
	sortBy := c.Query("sortBy")
	orderStr := c.Query("order")
	ascending := orderStr != "desc"

	// 转换排序字段
	var sortField utils.SortBy
	switch sortBy {
	case "name":
		sortField = utils.SortByName
	case "size":
		sortField = utils.SortBySize
	case "modTime":
		sortField = utils.SortByModTime
	default:
		sortField = utils.SortByName // 默认按名称排序
	}

	files, err := h.GetFileList(path, sortField, ascending)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取文件列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, files)
}

// HandleRenameFile 处理重命名文件的请求
func (h *FileController) HandleRenameFile(c *gin.Context) {
	var req struct {
		OldPath string `json:"oldPath"`
		NewPath string `json:"newPath"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}

	// 检查新路径是否已存在
	if _, err := os.Stat(req.NewPath); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "目标文件已存在"})
		return
	}

	// 执行重命名
	if err := os.Rename(req.OldPath, req.NewPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "重命名失败: " + err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// HandleBatchDeleteFiles 处理批量删除文件的请求
func (h *FileController) HandleBatchDeleteFiles(c *gin.Context) {
	var req struct {
		Paths []string `json:"paths"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if len(req.Paths) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未提供要删除的文件路径"})
		return
	}

	// 批量删除文件
	errors := h.BatchDeleteFiles(req.Paths)
	if len(errors) > 0 {
		// 如果有错误，返回第一个错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors[0].Error()})
		return
	}

	c.Status(http.StatusOK)
}
