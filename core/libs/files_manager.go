package libs

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

type FilesManager struct {
	files []FileInfo
}

func NewFilesManager() *FilesManager {
	return &FilesManager{}
}

// FileInfo 表示文件或目录的信息
type FileInfo struct {
	Name    string `json:"name"`    // 文件名
	Path    string `json:"path"`    // 完整路径
	Size    int64  `json:"size"`    // 文件大小（字节）
	IsDir   bool   `json:"isDir"`   // 是否是目录
	Mode    string `json:"mode"`    // 文件权限
	ModTime string `json:"modTime"` // 修改时间
	Owner   string `json:"owner"`   // Add owner field
	Group   string `json:"group"`   // Add group field
}

// GetFileList 获取指定目录下的文件列表
func (p *FilesManager) GetFileList(dirPath string) ([]FileInfo, error) {
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

		// Get file owner and group
		sys := info.Sys()
		var owner, group string
		if stat, ok := sys.(*syscall.Stat_t); ok {
			owner = strconv.FormatUint(uint64(stat.Uid), 10)
			group = strconv.FormatUint(uint64(stat.Gid), 10)
		}

		files = append(files, FileInfo{
			Name:    entry.Name(),
			Path:    fullPath,
			Size:    info.Size(),
			IsDir:   entry.IsDir(),
			Mode:    info.Mode().String(),
			ModTime: info.ModTime().Format("2006-01-02T15:04:05Z07:00"),
			Owner:   owner,
			Group:   group,
		})
	}

	return files, nil
}

// FormatFileSize converts bytes to human readable format
func FormatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
