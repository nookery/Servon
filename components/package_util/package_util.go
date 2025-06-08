// Package package_util 提供包管理功能
package package_util

import (
	"encoding/json"
	"os"
)

// ReadPackageVersion 尝试从 package.json 读取版本号
func ReadPackageVersion() string {
	data, err := os.ReadFile("package.json")
	if err != nil {
		return ""
	}

	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return ""
	}

	return pkg.Version
}

func IsDirExists(path string) bool {
	return os.MkdirAll(path, 0755) == nil
}

func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func GetDirSize(path string) int64 {
	files, err := os.ReadDir(path)
	if err != nil {
		return 0
	}

	var size int64
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		size += info.Size()
	}
	return size
}