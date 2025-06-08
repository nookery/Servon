// Package dev_util 提供开发环境检测功能
package dev_util

import (
	"os"
	"strings"
)

var DefaultDevUtil = &DevUtil{}

type DevUtil struct{}

// IsDev 判断是否为开发环境
func (d *DevUtil) IsDev() bool {
	if os.Args[0] == "main" ||
		strings.Contains(os.Args[0], "go-build") ||
		strings.Contains(os.Args[0], "/temp/servon") { // air 默认输出到 ./temp/servon
		return true
	}

	return false
}