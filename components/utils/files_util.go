package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/otiai10/copy"
)

var DefaultFileUtil = newFileUtil()

type FileUtil struct {
}

func newFileUtil() *FileUtil {
	return &FileUtil{}
}

// FileInfo 表示文件或目录的信息
type FileInfo struct {
	Name       string `json:"name"`       // 文件名
	Path       string `json:"path"`       // 完整路径
	Size       int64  `json:"size"`       // 文件大小（字节）
	IsDir      bool   `json:"isDir"`      // 是否是目录
	Mode       string `json:"mode"`       // 文件权限
	ModTime    string `json:"modTime"`    // 修改时间
	Owner      string `json:"owner"`      // 所有者
	Group      string `json:"group"`      // 组
	IsSymlink  bool   `json:"isSymlink"`  // 是否是软链接
	LinkTarget string `json:"linkTarget"` // 软链接目标路径
}

// SortBy 定义排序字段
type SortBy string

const (
	SortByName    SortBy = "name"
	SortBySize    SortBy = "size"
	SortByModTime SortBy = "modTime"
)

// GetFileList 获取指定目录下的文件列表
func (p *FileUtil) GetFileList(dirPath string, sortBy SortBy, ascending bool) ([]FileInfo, error) {
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
		// if strings.HasPrefix(entry.Name(), ".") {
		// 	continue
		// }

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

		// 检查是否是软链接
		isSymlink := false
		var linkTarget string
		fileInfo, err := os.Lstat(fullPath)
		if err == nil && fileInfo.Mode()&os.ModeSymlink != 0 {
			isSymlink = true
			// 获取软链接目标路径
			if target, err := os.Readlink(fullPath); err == nil {
				linkTarget = target
			}
		}

		files = append(files, FileInfo{
			Name:       entry.Name(),
			Path:       fullPath,
			Size:       info.Size(),
			IsDir:      entry.IsDir(),
			Mode:       info.Mode().String(),
			ModTime:    info.ModTime().Format("2006-01-02T15:04:05Z07:00"),
			Owner:      owner,
			Group:      group,
			IsSymlink:  isSymlink,
			LinkTarget: linkTarget,
		})
	}

	// 根据指定字段排序
	switch sortBy {
	case SortByName:
		if ascending {
			sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })
		} else {
			sort.Slice(files, func(i, j int) bool { return files[i].Name > files[j].Name })
		}
	case SortBySize:
		if ascending {
			sort.Slice(files, func(i, j int) bool { return files[i].Size < files[j].Size })
		} else {
			sort.Slice(files, func(i, j int) bool { return files[i].Size > files[j].Size })
		}
	case SortByModTime:
		if ascending {
			sort.Slice(files, func(i, j int) bool { return files[i].ModTime < files[j].ModTime })
		} else {
			sort.Slice(files, func(i, j int) bool { return files[i].ModTime > files[j].ModTime })
		}
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

// IsDirExists 判断目录是否存在
func (p *FileUtil) IsDirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// IsFileExists 判断文件是否存在
func (p *FileUtil) IsFileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// IsExist 判断文件或目录是否存在
func (p *FileUtil) IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetDirSize 获取目录大小
func (p *FileUtil) GetDirSize(dirPath string) int64 {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return 0
	}
	return info.Size()
}

// CopyDir 复制目录
func (p *FileUtil) CopyDir(srcDir, destDir string) error {
	// 确保源目录存在
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		return fmt.Errorf("源目录不存在或无法访问: %v", err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("源路径不是目录: %s", srcDir)
	}

	// 使用 otiai10/copy 库复制目录
	opt := copy.Options{
		Skip: func(srcinfo os.FileInfo, src string, dest string) (bool, error) {
			return false, nil // 不跳过任何文件
		},
		OnSymlink: func(src string) copy.SymlinkAction {
			return copy.Deep // 复制软链接指向的内容
		},
		OnDirExists: func(src, dest string) copy.DirExistsAction {
			return copy.Merge // 合并目录
		},
	}

	return copy.Copy(srcDir, destDir, opt)
}

// CopyFile 复制文件
func (p *FileUtil) CopyFile(srcPath, destPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	return err
}

// ListFiles 列出目录下的文件
func (p *FileUtil) ListFiles(dirPath string) ([]FileInfo, error) {
	files, err := p.GetFileList(dirPath, SortByName, true)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// ListFilesWithLimit 列出目录下的文件，并限制数量
func (p *FileUtil) ListFilesWithLimit(dirPath string, limit int) ([]FileInfo, error) {
	files, err := p.ListFiles(dirPath)
	if err != nil {
		return nil, err
	}

	return files[:limit], nil
}

// SymlinkForce 强制创建软链接
func (p *FileUtil) SymlinkForce(oldname, newname string) error {
	// 如果目标路径存在，则删除
	if _, err := os.Stat(newname); err == nil {
		err = os.Remove(newname)
		if err != nil {
			return fmt.Errorf("删除目标路径失败: %v", err)
		}
	}

	// 检查目标路径是否成功删除
	if _, err := os.Stat(newname); err == nil {
		return fmt.Errorf("目标路径删除失败: %v", err)
	}

	// 创建软链接
	return p.Symlink(oldname, newname)
}

// Symlink 创建软链接
func (p *FileUtil) Symlink(oldname, newname string) error {
	// 获取目标路径的目录
	targetDir := filepath.Dir(newname)

	// 检查目标目录是否存在，如果不存在则创建
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("创建目标目录失败: %v", err)
		}
	}

	// 创建符号链接
	err := os.Symlink(oldname, newname)
	if err != nil {
		return fmt.Errorf("os.Symlink 创建软链接失败: %v。oldname: %s, newname: %s", err, oldname, newname)
	}

	return nil
}

// MakeDir 创建目录
func (p *FileUtil) MakeDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// RemoveDir 删除目录
func (p *FileUtil) RemoveDir(dirPath string) error {
	// 如果不存在，则忽略
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil
	}

	return os.RemoveAll(dirPath)
}

// RemoveFile 删除文件
func (p *FileUtil) RemoveFile(filePath string) error {
	// 如果不存在，则忽略
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(filePath)
}

// RemoveFileOrDir 删除文件或目录
func (p *FileUtil) RemoveFileOrDir(path string) error {
	// 如果不存在，则忽略
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	return os.RemoveAll(path)
}

// 检查文件是否是软链接
func (p *FileUtil) IsSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}
