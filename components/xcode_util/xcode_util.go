package xcode_util

//
// 这个组件封装了与Xcode开发环境相关的工具函数，
// 包括构建环境信息显示、命令执行等功能。

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gookit/color"
)

var DefaultXcodeUtil = &XcodeUtil{}
var ShowBuildEnvironmentInfo = DefaultXcodeUtil.ShowBuildEnvironmentInfo
var GetCommandOutput = DefaultXcodeUtil.GetCommandOutput

type XcodeUtil struct{}

// NewXcodeUtil 创建新的XcodeUtil实例
func NewXcodeUtil() *XcodeUtil {
	return &XcodeUtil{}
}

// ShowBuildEnvironmentInfo 显示构建环境信息
func (x *XcodeUtil) ShowBuildEnvironmentInfo(scheme, buildPath, arch string, verbose bool) {
	color.Blue.Println("===========================================")
	color.Blue.Println("         应用构建环境信息                ")
	color.Blue.Println("===========================================")
	fmt.Println()
	
	// 系统信息
	color.Green.Println("📱 系统信息:")
	fmt.Printf("   操作系统: %s %s\n", runtime.GOOS, runtime.GOARCH)
	if hostname, err := os.Hostname(); err == nil {
		fmt.Printf("   主机名称: %s\n", hostname)
	}
	fmt.Println()
	
	// Xcode 信息
	color.Green.Println("🔨 Xcode 开发环境:")
	if xcodeVersion := x.GetCommandOutput("xcodebuild", "-version"); xcodeVersion != "" {
		lines := strings.Split(xcodeVersion, "\n")
		if len(lines) >= 1 {
			fmt.Printf("   Xcode 版本: %s\n", lines[0])
		}
		if len(lines) >= 2 {
			fmt.Printf("   构建版本: %s\n", lines[1])
		}
	}
	if sdkPath := x.GetCommandOutput("xcrun", "--show-sdk-path"); sdkPath != "" {
		fmt.Printf("   SDK 路径: %s\n", sdkPath)
	}
	if devDir := x.GetCommandOutput("xcode-select", "-p"); devDir != "" {
		fmt.Printf("   开发者目录: %s\n", devDir)
	}
	fmt.Println()
	
	// Swift 信息
	color.Green.Println("🚀 Swift 编译器:")
	if swiftVersion := x.GetCommandOutput("swift", "--version"); swiftVersion != "" {
		lines := strings.Split(swiftVersion, "\n")
		if len(lines) >= 1 {
			fmt.Printf("   Swift 版本: %s\n", lines[0])
		}
	}
	fmt.Println()
	
	// Git 信息
	color.Green.Println("📝 Git 版本控制:")
	if gitVersion := x.GetCommandOutput("git", "--version"); gitVersion != "" {
		fmt.Printf("   Git 版本: %s\n", gitVersion)
	}
	if branch := x.GetCommandOutput("git", "branch", "--show-current"); branch != "" {
		fmt.Printf("   当前分支: %s\n", branch)
	}
	if commit := x.GetCommandOutput("git", "log", "-1", "--pretty=format:%h - %s (%an, %ar)"); commit != "" {
		fmt.Printf("   最新提交: %s\n", commit)
	}
	fmt.Println()
	
	// 构建环境变量
	color.Green.Println("🌍 构建环境变量:")
	fmt.Printf("   构建方案: %s\n", scheme)
	fmt.Printf("   构建路径: %s\n", buildPath)
	fmt.Printf("   目标架构: %s\n", arch)
	fmt.Printf("   构建配置: Release\n")
	fmt.Printf("   详细日志: %t\n", verbose)
	if cwd, err := os.Getwd(); err == nil {
		fmt.Printf("   工作目录: %s\n", cwd)
	}
	fmt.Println()
}

// GetCommandOutput 执行命令并返回输出
func (x *XcodeUtil) GetCommandOutput(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}