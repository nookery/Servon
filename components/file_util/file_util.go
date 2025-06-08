// Package file_util 提供文件和目录操作功能
//
// 这个组件封装了与文件系统相关的工具函数，
// 包括文件读写、目录操作、文件复制等功能。
package file_util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
)

var DefaultFileUtil = newFileUtil()

type FileUtil struct {
}

func newFileUtil() *FileUtil {
	return &FileUtil{}
}

func NewFileUtil() *FileUtil {
	return &FileUtil{}
}

// Exists 检查文件或目录是否存在
func (f *FileUtil) Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDir 检查路径是否为目录
func (f *FileUtil) IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsFile 检查路径是否为文件
func (f *FileUtil) IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// CreateDir 创建目录
func (f *FileUtil) CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// RemoveDir 删除目录
func (f *FileUtil) RemoveDir(path string) error {
	return os.RemoveAll(path)
}

// RemoveFile 删除文件
func (f *FileUtil) RemoveFile(path string) error {
	return os.Remove(path)
}

// CopyFile 复制文件
func (f *FileUtil) CopyFile(src, dst string) error {
	return copy.Copy(src, dst)
}

// MoveFile 移动文件
func (f *FileUtil) MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}

// ReadFile 读取文件内容
func (f *FileUtil) ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WriteFile 写入文件内容
func (f *FileUtil) WriteFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// AppendFile 追加文件内容
func (f *FileUtil) AppendFile(path, content string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	return err
}

// ListDir 列出目录内容
func (f *FileUtil) ListDir(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names, nil
}

// GetFileSize 获取文件大小
func (f *FileUtil) GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// GetFileInfo 获取文件信息
func (f *FileUtil) GetFileInfo(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// GetAbsPath 获取绝对路径
func (f *FileUtil) GetAbsPath(path string) (string, error) {
	return filepath.Abs(path)
}

// GetRelPath 获取相对路径
func (f *FileUtil) GetRelPath(basepath, targpath string) (string, error) {
	return filepath.Rel(basepath, targpath)
}

// JoinPath 连接路径
func (f *FileUtil) JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

// GetDir 获取文件所在目录
func (f *FileUtil) GetDir(path string) string {
	return filepath.Dir(path)
}

// GetBaseName 获取文件名（不含路径）
func (f *FileUtil) GetBaseName(path string) string {
	return filepath.Base(path)
}

// GetExt 获取文件扩展名
func (f *FileUtil) GetExt(path string) string {
	return filepath.Ext(path)
}

// ChangeExt 更改文件扩展名
func (f *FileUtil) ChangeExt(path, newExt string) string {
	dir := filepath.Dir(path)
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return filepath.Join(dir, base+newExt)
}

// FindFiles 查找文件
func (f *FileUtil) FindFiles(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			matches = append(matches, path)
		}
		return nil
	})
	return matches, err
}

// GetTempDir 获取临时目录
func (f *FileUtil) GetTempDir() string {
	return os.TempDir()
}

// CreateTempFile 创建临时文件
func (f *FileUtil) CreateTempFile(pattern string) (*os.File, error) {
	return os.CreateTemp("", pattern)
}

// CreateTempDir 创建临时目录
func (f *FileUtil) CreateTempDir(pattern string) (string, error) {
	return os.MkdirTemp("", pattern)
}

// IsEmpty 检查目录是否为空
func (f *FileUtil) IsEmpty(path string) (bool, error) {
	if !f.IsDir(path) {
		return false, fmt.Errorf("%s is not a directory", path)
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

// GetDirSize 获取目录大小
func (f *FileUtil) GetDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// FormatFileSize 格式化文件大小
func (f *FileUtil) FormatFileSize(size int64) string {
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

// SetPermissions 设置文件权限
func (f *FileUtil) SetPermissions(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

// GetPermissions 获取文件权限
func (f *FileUtil) GetPermissions(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Mode(), nil
}
