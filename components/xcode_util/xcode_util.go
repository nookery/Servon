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

// IsXcodeProjectDirectory 检查指定目录是否为Xcode项目目录
// 检查是否存在 .xcodeproj 或 .xcworkspace 文件
func (x *XcodeUtil) IsXcodeProjectDirectory(dir string) (bool, string, error) {
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return false, "", fmt.Errorf("无法获取当前目录: %v", err)
		}
	}

	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false, "", fmt.Errorf("目录不存在: %s", dir)
	}

	// 读取目录内容
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, "", fmt.Errorf("无法读取目录: %v", err)
	}

	// 检查是否存在 .xcworkspace 文件（优先级更高）
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".xcworkspace") {
			return true, "workspace", nil
		}
	}

	// 检查是否存在 .xcodeproj 文件
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".xcodeproj") {
			return true, "project", nil
		}
	}

	return false, "", nil
}

// DetectProjectFile 检测指定目录中的Xcode项目文件
// 参数:
//   - dir: 要搜索的目录，如果为空则使用当前目录
//   - verbose: 是否输出详细信息
//
// 返回值:
//   - 项目文件路径
//   - 项目类型 ("workspace" 或 "project")
//   - 错误信息
func (x *XcodeUtil) DetectProjectFile(dir string, verbose bool) (string, string, error) {
	var searchDir string
	if dir != "" {
		// 如果指定了目录，使用指定的目录
		searchDir = dir
		if verbose {
			color.Info.Printf("🔍 在指定目录中搜索: %s\n", dir)
		}
	} else {
		// 否则使用当前目录
		cwd, err := os.Getwd()
		if err != nil {
			return "", "", fmt.Errorf("无法获取当前目录: %v", err)
		}
		searchDir = cwd
		if verbose {
			color.Info.Printf("🔍 在当前目录中搜索: %s\n", cwd)
		}
	}

	// 读取目录内容
	entries, err := os.ReadDir(searchDir)
	if err != nil {
		return "", "", fmt.Errorf("无法读取目录 %s: %v", searchDir, err)
	}

	// 查找 .xcworkspace 文件（优先级更高）
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".xcworkspace") {
			workspacePath := fmt.Sprintf("%s/%s", searchDir, entry.Name())
			if verbose {
				color.Success.Printf("✅ 找到 workspace 文件: %s\n", workspacePath)
			}
			return workspacePath, "workspace", nil
		}
	}

	// 查找 .xcodeproj 文件
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".xcodeproj") {
			projectPath := fmt.Sprintf("%s/%s", searchDir, entry.Name())
			if verbose {
				color.Success.Printf("✅ 找到 project 文件: %s\n", projectPath)
			}
			return projectPath, "project", nil
		}
	}

	return "", "", fmt.Errorf("在目录 %s 中未找到 .xcodeproj 或 .xcworkspace 文件", searchDir)
}

// DetectScheme 自动检测指定目录中可用的第一个 scheme
// 参数:
//   - dir: 要搜索的目录，如果为空则使用当前目录
//
// 返回值:
//   - 检测到的第一个 scheme 名称，如果未找到则返回空字符串
func (x *XcodeUtil) DetectScheme(dir string) string {
	projectFile, projectType, err := x.DetectProjectFile(dir, false)
	if err != nil {
		return ""
	}

	var cmd *exec.Cmd
	if projectType == "workspace" {
		cmd = exec.Command("xcodebuild", "-workspace", projectFile, "-list")
	} else {
		cmd = exec.Command("xcodebuild", "-project", projectFile, "-list")
	}

	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// 解析输出，查找 schemes
	lines := strings.Split(string(output), "\n")
	inSchemes := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "Schemes:" {
			inSchemes = true
			continue
		}
		if inSchemes && line != "" && !strings.Contains(line, ":") {
			return line // 返回第一个找到的 scheme
		}
		if inSchemes && line == "" {
			break
		}
	}

	return ""
}

// GetBuildTargetInfo 获取构建目标信息
// 参数:
//   - projectFile: 项目文件路径
//   - projectType: 项目类型 ("workspace" 或 "project")
//   - scheme: 构建方案名称
//   - arch: 目标架构
//
// 返回值:
//   - BuildTargetInfo 结构体包含所有构建目标信息
//   - 错误信息
func (x *XcodeUtil) GetBuildTargetInfo(projectFile, projectType, scheme, arch string) (*BuildTargetInfo, error) {
	info := &BuildTargetInfo{
		ProjectFile: projectFile,
		ProjectType: projectType,
		Scheme:      scheme,
		TargetArch:  arch,
	}

	// 设置项目类型显示名称
	if projectType == "workspace" {
		info.ProjectTypeName = "Xcode Workspace"
	} else {
		info.ProjectTypeName = "Xcode Project"
	}

	// 获取项目支持的架构
	var cmd *exec.Cmd
	if projectType == "workspace" {
		cmd = exec.Command("xcodebuild", "-workspace", projectFile, "-scheme", scheme, "-showBuildSettings", "-configuration", "Release")
	} else {
		cmd = exec.Command("xcodebuild", "-project", projectFile, "-scheme", scheme, "-showBuildSettings", "-configuration", "Release")
	}

	output, err := cmd.Output()
	if err != nil {
		// 如果无法获取架构信息，不返回错误，只是留空
		info.ProjectArchs = "未知"
	} else {
		// 解析架构信息
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "ARCHS =") {
				parts := strings.Split(line, "=")
				if len(parts) >= 2 {
					info.ProjectArchs = strings.TrimSpace(parts[1])
					break
				}
			}
		}
		if info.ProjectArchs == "" {
			info.ProjectArchs = "未知"
		}
	}

	return info, nil
}
