package managers

import (
	"servon/core/internal/utils"
)

type FileManager struct {
	*utils.FileUtil
}

func NewFileManager() *FileManager {
	return &FileManager{FileUtil: utils.DefaultFileUtil}
}

// GetFileList 获取文件列表，支持排序
func (m *FileManager) GetFileList(path string, sortBy utils.SortBy, ascending bool) ([]utils.FileInfo, error) {
	return m.FileUtil.GetFileList(path, sortBy, ascending)
}
