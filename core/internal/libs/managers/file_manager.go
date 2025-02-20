package managers

import (
	"servon/core/internal/libs/utils"
)

type FileManager struct {
	*utils.FileUtil
}

func NewFileManager() *FileManager {
	return &FileManager{FileUtil: utils.DefaultFileUtil}
}
