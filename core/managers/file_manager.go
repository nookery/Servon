package managers

import (
	"fmt"
	"os"
	"servon/core/utils"
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

// BatchDeleteFiles 批量删除文件，返回错误列表
func (m *FileManager) BatchDeleteFiles(paths []string) []error {
	var errors []error
	for _, path := range paths {
		if err := os.Remove(path); err != nil {
			errors = append(errors, fmt.Errorf("删除文件 %s 失败: %v", path, err))
		}
	}
	return errors
}

// DeleteFile 删除单个文件
func (m *FileManager) DeleteFile(path string) error {
	return os.Remove(path)
}
