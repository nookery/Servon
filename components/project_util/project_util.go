// Package project_util 提供项目检测和管理功能
package project_util

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

	// Vue 项目特征
	if fileExists(filepath.Join(projectPath, "vue.config.js")) ||
		fileExists(filepath.Join(projectPath, "vite.config.js")) ||
		fileExists(filepath.Join(projectPath, "vite.config.ts")) {
		return "vue"
	}

	// React 项目特征
	if fileExists(filepath.Join(projectPath, "package.json")) {
		// 可以进一步检查 package.json 内容来确定是否为 React 项目
		return "react"
	}

	// Node.js 项目特征
	if fileExists(filepath.Join(projectPath, "package.json")) {
		return "nodejs"
	}

	// Go 项目特征
	if fileExists(filepath.Join(projectPath, "go.mod")) {
		return "go"
	}

	// Python 项目特征
	if fileExists(filepath.Join(projectPath, "requirements.txt")) ||
		fileExists(filepath.Join(projectPath, "setup.py")) ||
		fileExists(filepath.Join(projectPath, "pyproject.toml")) {
		return "python"
	}

	// Java 项目特征
	if fileExists(filepath.Join(projectPath, "pom.xml")) ||
		fileExists(filepath.Join(projectPath, "build.gradle")) {
		return "java"
	}

	return "unknown"
}

// fileExists 检查文件是否存在
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// GetProjectRoot 获取项目根目录
func (p *ProjectUtil) GetProjectRoot(startPath string) string {
	currentPath := startPath
	for {
		// 检查常见的项目根目录标识文件
		if fileExists(filepath.Join(currentPath, ".git")) ||
			fileExists(filepath.Join(currentPath, "go.mod")) ||
			fileExists(filepath.Join(currentPath, "package.json")) ||
			fileExists(filepath.Join(currentPath, "composer.json")) {
			return currentPath
		}

		// 向上一级目录
		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			// 已经到达根目录
			break
		}
		currentPath = parentPath
	}

	return startPath
}