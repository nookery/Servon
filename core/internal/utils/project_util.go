package utils

import (
	"os"
	"path/filepath"
)

var DefaultProjectUtil = NewProjectUtil()

type ProjectUtil struct {
}

func NewProjectUtil() *ProjectUtil {
	return &ProjectUtil{}
}

// DetectProjectType 通过检查项目目录特征文件来识别项目类型
func (p *ProjectUtil) DetectProjectType(projectPath string) string {
	// Laravel 项目特征
	if fileExists(filepath.Join(projectPath, "artisan")) &&
		fileExists(filepath.Join(projectPath, "composer.json")) {
		return "laravel"
	}

	// Astro 项目特征
	if fileExists(filepath.Join(projectPath, "astro.config.mjs")) ||
		fileExists(filepath.Join(projectPath, "astro.config.js")) {
		return "astro"
	}

	// Flutter 项目特征
	if fileExists(filepath.Join(projectPath, "pubspec.yaml")) &&
		dirExists(filepath.Join(projectPath, "lib")) {
		return "flutter"
	}

	return "unknown"
}

// fileExists 检查文件是否存在
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// dirExists 检查目录是否存在
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

//
