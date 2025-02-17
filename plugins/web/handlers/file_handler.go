package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// HandleDeleteFile 处理删除文件的请求
func (h *WebHandler) HandleDeleteFile(c *gin.Context) {
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
func (h *WebHandler) HandleCreateFile(c *gin.Context) {
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
func (h *WebHandler) HandleFileContent(c *gin.Context) {
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
func (h *WebHandler) HandleSaveFile(c *gin.Context) {
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
func (h *WebHandler) HandleFileDownload(c *gin.Context) {
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
func (h *WebHandler) HandleFileList(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		path = "/"
	}

	files, err := h.App.FileUtil.GetFileList(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取文件列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, files)
}

// HandleRenameFile 处理重命名文件的请求
func (h *WebHandler) HandleRenameFile(c *gin.Context) {
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
