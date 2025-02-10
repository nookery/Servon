package libs

import (
	"os"
	"path/filepath"
	"strings"
)

// FileInfo 表示文件或目录的信息
type FileInfo struct {
	Name    string `json:"name"`    // 文件名
	Path    string `json:"path"`    // 完整路径
	Size    int64  `json:"size"`    // 文件大小（字节）
	IsDir   bool   `json:"isDir"`   // 是否是目录
	Mode    string `json:"mode"`    // 文件权限
	ModTime string `json:"modTime"` // 修改时间
}

// GetFileList 获取指定目录下的文件列表
func GetFileList(dirPath string) ([]FileInfo, error) {
	// 清理和规范化路径
	dirPath = filepath.Clean(dirPath)
	if !strings.HasPrefix(dirPath, "/") {
		dirPath = "/" + dirPath
	}

	// 打开目录
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// 读取目录内容
	entries, err := dir.ReadDir(-1)
	if err != nil {
		return nil, err
	}

	// 构建文件信息列表
	files := make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		// 跳过隐藏文件
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// 获取详细信息
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// 构建完整路径
		fullPath := filepath.Join(dirPath, entry.Name())

		files = append(files, FileInfo{
			Name:    entry.Name(),
			Path:    fullPath,
			Size:    info.Size(),
			IsDir:   entry.IsDir(),
			Mode:    info.Mode().String(),
			ModTime: info.ModTime().Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return files, nil
}
