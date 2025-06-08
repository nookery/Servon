package xcode

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"servon/components/xcode_util"
	"strings"

	"github.com/gookit/color"
)

// showGitStatus 显示 Git 状态
func showGitStatus() {
	color.Green.Println("📝 Git 状态变更:")

	if status := xcode_util.DefaultXcodeUtil.GetCommandOutput("git", "status", "--porcelain"); status != "" {
		lines := strings.Split(strings.TrimSpace(status), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				fmt.Printf("   %s\n", line)
			}
		}
	} else {
		fmt.Println("   无变更")
	}
	fmt.Println()
}

// findPbxprojFile 自动查找 .pbxproj 文件
func findPbxprojFile() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("无法获取当前目录: %v", err)
	}

	// 在当前目录及其子目录中查找 .pbxproj 文件（排除 Resources 和 temp 目录）
	var projectFile string
	err = filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 忽略错误，继续查找
		}

		// 跳过深度超过2层的目录
		relPath, _ := filepath.Rel(cwd, path)
		if strings.Count(relPath, string(filepath.Separator)) > 2 {
			return filepath.SkipDir
		}

		// 跳过 Resources 和 temp 目录
		if info.IsDir() && (strings.Contains(path, "Resources") || strings.Contains(path, "temp")) {
			return filepath.SkipDir
		}

		// 查找 .pbxproj 文件
		if strings.HasSuffix(path, ".pbxproj") {
			projectFile = path
			return fmt.Errorf("found") // 用错误来停止遍历
		}

		return nil
	})

	if projectFile == "" {
		return "", fmt.Errorf("未找到 .pbxproj 配置文件")
	}

	return projectFile, nil
}

// showBumpEnvironmentInfo 显示版本管理环境信息
func showBumpEnvironmentInfo() {
	color.Blue.Println("===========================================")
	color.Blue.Println("         版本管理环境信息                ")
	color.Blue.Println("===========================================")
	fmt.Println()

	// 系统信息
	color.Green.Println("📱 系统信息:")
	if hostname, err := os.Hostname(); err == nil {
		fmt.Printf("   主机名称: %s\n", hostname)
	}
	if cwd, err := os.Getwd(); err == nil {
		fmt.Printf("   当前目录: %s\n", cwd)
	}
	fmt.Println()

	// Git 信息
	color.Green.Println("📝 Git 版本控制:")
	if gitVersion := xcode_util.DefaultXcodeUtil.GetCommandOutput("git", "--version"); gitVersion != "" {
		fmt.Printf("   Git 版本: %s\n", gitVersion)
	}
	if branch := xcode_util.DefaultXcodeUtil.GetCommandOutput("git", "branch", "--show-current"); branch != "" {
		fmt.Printf("   当前分支: %s\n", branch)
	}
	if commit := xcode_util.DefaultXcodeUtil.GetCommandOutput("git", "log", "-1", "--pretty=format:%h - %s (%an, %ar)"); commit != "" {
		fmt.Printf("   最新提交: %s\n", commit)
	}
	fmt.Println()
}

// showDevelopmentRoadmap 显示开发路线图
func showDevelopmentRoadmap(currentStep string) {
	fmt.Println()
	color.Blue.Println("===========================================")
	color.Blue.Println("         🗺️  开发分发路线图                ")
	color.Blue.Println("===========================================")
	fmt.Println()

	steps := []string{
		"setup:⚙️ 环境设置:配置代码签名环境",
		"version:📝 版本管理:查看或更新应用版本号",
		"build:🔨 构建应用:编译源代码，生成可执行文件",
		"codesign:🔐 代码签名:为应用添加数字签名，确保安全性",
		"package:📦 打包分发:创建 DMG 安装包",
		"notarize:✅ 公证验证:Apple 官方验证（可选）",
		"distribute:🚀 发布分发:上传到分发平台或直接分发",
	}

	color.Cyan.Print("📍 当前位置: ")
	switch currentStep {
	case "setup":
		color.Green.Println("环境设置")
	case "version":
		color.Green.Println("版本管理")
	case "build":
		color.Green.Println("构建应用")
	case "codesign":
		color.Green.Println("代码签名")
	case "package":
		color.Green.Println("打包分发")
	case "notarize":
		color.Green.Println("公证验证")
	case "distribute":
		color.Green.Println("发布分发")
	default:
		color.Yellow.Println("未知步骤")
	}
	fmt.Println()

	// 显示路线图
	for _, step := range steps {
		parts := strings.Split(step, ":")
		stepId := parts[0]
		stepIcon := parts[1]
		stepDesc := parts[2]

		if stepId == currentStep {
			color.Green.Printf("▶ %s %s\n", stepIcon, stepDesc)
		} else {
			fmt.Printf("  %s %s\n", stepIcon, stepDesc)
		}
	}

	fmt.Println()
	color.Yellow.Println("💡 下一步建议:")
	switch currentStep {
	case "setup":
		color.Cyan.Println("   查看版本信息: go run main.go xcode version")
		color.Cyan.Println("   或直接构建应用: go run main.go xcode build")
	case "version":
		color.Cyan.Println("   构建应用: go run main.go xcode build")
	case "build":
		color.Cyan.Println("   运行代码签名: go run main.go xcode codesign")
	case "codesign":
		color.Cyan.Println("   创建安装包: go run main.go xcode package")
	case "package":
		fmt.Println("   进行公证验证或直接分发应用")
	case "notarize":
		fmt.Println("   发布到分发平台或提供下载链接")
	case "distribute":
		fmt.Println("   🎉 开发分发流程已完成！")
	}

	fmt.Println()
	color.Blue.Println("===========================================")
}

// showAvailableSchemes 显示可用的 schemes
//
//   - getAvailableSchemes 获取项目中可用的 scheme 列表
//   - workDir: 工作目录，如果为空则使用当前目录
//   - verbose: 是否输出详细信息
//
// 返回值:
//   - schemes列表
//   - 项目文件路径
//   - 项目类型
//   - 错误信息
func getAvailableSchemes(workDir string, verbose bool) ([]string, string, string, error) {
	if verbose {
		color.Yellow.Println("正在检查项目中可用的 scheme...")
	}

	projectFile, projectType, err := detectProjectFile(workDir, verbose)
	if err != nil {
		return nil, "", "", err
	}

	if verbose {
		color.Green.Printf("在项目 %s 中找到以下可用的 scheme:\n", projectFile)
	}

	var cmd *exec.Cmd
	if projectType == "workspace" {
		cmd = exec.Command("xcodebuild", "-workspace", projectFile, "-list")
	} else {
		cmd = exec.Command("xcodebuild", "-project", projectFile, "-list")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, projectFile, projectType, fmt.Errorf("无法获取 scheme 列表: %v", err)
	}

	// 解析 schemes
	lines := strings.Split(string(output), "\n")
	inSchemes := false
	var schemes []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "Schemes:" {
			inSchemes = true
			continue
		}
		if inSchemes && line != "" && !strings.Contains(line, ":") {
			schemes = append(schemes, line)
		}
		if inSchemes && line == "" {
			break
		}
	}

	return schemes, projectFile, projectType, nil
}

// showAvailableSchemes 显示项目中可用的 scheme 列表（保持向后兼容）
func showAvailableSchemes(workDir string, verbose bool) {
	schemes, projectFile, _, err := getAvailableSchemes(workDir, verbose)
	if err != nil {
		color.Error.Printf("❌ %s\n", err.Error())
		return
	}

	if verbose {
		color.Green.Printf("在项目 %s 中找到以下可用的 scheme:\n", projectFile)
	}

	// 显示 schemes
	for _, scheme := range schemes {
		fmt.Printf("  - %s\n", scheme)
	}

	fmt.Println()
	color.Yellow.Println("💡 使用示例:")
	color.Cyan.Println("   go run main.go xcode build --scheme YourSchemeName")
}

// detectProjectFile 检测项目文件
// detectProjectFile 检测项目文件（已废弃，使用 xcode_util.DetectProjectFile）
func detectProjectFile(workDir string, verbose bool) (string, string, error) {
	return xcode_util.DefaultXcodeUtil.DetectProjectFile(workDir, verbose)
}

// detectScheme 自动检测可用的 scheme（已废弃，使用 xcode_util.DetectScheme）
func detectScheme() string {
	return xcode_util.DefaultXcodeUtil.DetectScheme("")
}
