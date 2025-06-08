package xcode

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"servon/components/xcode_util"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "构建 iOS/macOS 应用",
	Long:  color.Success.Render("\r\n构建 iOS/macOS 应用程序，支持多架构和详细日志输出"),
	Run: func(cmd *cobra.Command, args []string) {
		scheme, _ := cmd.Flags().GetString("scheme")
		buildPath, _ := cmd.Flags().GetString("build-path")
		arch, _ := cmd.Flags().GetString("arch")
		verbose, _ := cmd.Flags().GetBool("verbose")
		clean, _ := cmd.Flags().GetBool("clean")
		workDir, _ := cmd.Flags().GetString("workdir")

		// 如果指定了工作目录，切换到该目录
		if workDir != "" {
			if err := os.Chdir(workDir); err != nil {
				color.Error.Printf("❌ 无法切换到工作目录 %s: %s\n", workDir, err.Error())
				os.Exit(1)
			}
			if verbose {
				color.Info.Printf("📁 切换到工作目录: %s\n", workDir)
			}
		}

		// 显示环境信息
		xcode_util.DefaultXcodeUtil.ShowBuildEnvironmentInfo(scheme, buildPath, arch, verbose)

		// 检查当前目录是否为Xcode项目目录
		isXcodeProject, projectType, err := xcode_util.DefaultXcodeUtil.IsXcodeProjectDirectory("")
		if err != nil {
			color.Error.Printf("❌ 检查项目目录时出错: %s\n", err.Error())
			os.Exit(1)
		}
		if !isXcodeProject {
			color.Warnln("❌ 错误: 当前目录不是有效的Xcode项目目录")
			color.Yellow.Println("💡 请确保当前目录包含 .xcodeproj 或 .xcworkspace 文件")
			os.Exit(0)
		}
		if verbose {
			color.Green.Printf("✅ 检测到有效的Xcode项目 (类型: %s)\n", projectType)
		}

		// 检查必需的环境变量
		if scheme == "" {
			color.Warnln("❌ 错误: 未指定构建方案 (scheme)")
			// 获取可用的 schemes
			schemes, projectFile, _, err := getAvailableSchemes("", verbose)
			if err == nil && len(schemes) > 0 {
				color.Green.Printf("📋 在项目 %s 中找到以下可用的构建方案:\n", projectFile)
				fmt.Println()
				for i, s := range schemes {
					fmt.Printf("  %d. %s\n", i+1, s)
				}
				fmt.Println()
				color.Yellow.Println("💡 请使用以下命令指定构建方案:")
				color.Cyan.Printf("   servon xcode build --scheme <方案名称>\n")
				fmt.Println()
				color.Yellow.Println("📝 示例:")
				for _, s := range schemes {
					color.Cyan.Printf("   servon xcode build --scheme %s\n", s)
				}
			} else {
				color.Error.Println("❌ 无法获取可用的构建方案")
				if err != nil {
					color.Error.Printf("   错误详情: %s\n", err.Error())
				}
			}
			os.Exit(0)
		}

		// 设置默认值
		if buildPath == "" {
			buildPath = "./temp"
		}

		// 检测项目文件
		projectFile, projectType, err := detectProjectFile("", verbose)
		if err != nil {
			color.Error.Printf("❌ %s\n", err.Error())
			os.Exit(1)
		}

		// 如果未指定架构，构建所有架构
		var architectures []string
		if arch == "" {
			architectures = []string{"x86_64", "arm64"}
			color.Info.Println("🏗️  未指定架构，将构建所有支持的架构: x86_64, arm64")
		} else {
			architectures = []string{arch}
		}

		// 为每个架构执行构建
		var buildPaths []string
		for _, currentArch := range architectures {
			// 根据架构设置构建路径
			archBuildPath := filepath.Join(buildPath, currentArch)
			buildPaths = append(buildPaths, archBuildPath)

			if verbose {
				color.Info.Printf("构建路径 (%s): %s\n", currentArch, archBuildPath)
			}

			// 显示构建目标信息
			showBuildTargetInfo(projectFile, projectType, scheme, currentArch)

			// 执行构建
			err = performBuild(projectFile, projectType, scheme, archBuildPath, currentArch, verbose, clean)
			if err != nil {
				color.Error.Printf("❌ 构建失败 (%s): %s\n", currentArch, err.Error())
				os.Exit(1)
			}

			color.Success.Printf("✅ %s 架构构建成功完成！\n", currentArch)
		}

		// 显示所有构建产物位置
		color.Success.Println("🎉 所有架构构建成功完成！")
		color.Green.Println("📦 构建产物位置:")
		for i, buildPath := range buildPaths {
			color.Green.Printf("   %s: %s/Build/Products/Release/\n", architectures[i], buildPath)
		}

		// 显示开发路线图
		showDevelopmentRoadmap("build")
	},
}

func init() {
	buildCmd.Flags().StringP("scheme", "s", "", "构建方案名称")
	buildCmd.Flags().StringP("build-path", "b", "./temp", "构建输出路径")
	buildCmd.Flags().StringP("arch", "a", "", "目标架构 (不指定则构建所有架构: x86_64, arm64; 可选: universal, x86_64, arm64)")
	buildCmd.Flags().BoolP("verbose", "v", false, "显示详细构建日志")
	buildCmd.Flags().Bool("clean", true, "构建前清理")
	buildCmd.Flags().StringP("workdir", "w", "", "指定工作目录，在该目录中搜索项目文件")
}

// showBuildTargetInfo 显示构建目标信息
// showBuildTargetInfo 显示构建目标信息
// @deprecated: Use xcode_util.DefaultXcodeUtil.GetBuildTargetInfo instead.
func showBuildTargetInfo(projectFile, projectType, scheme, arch string) {
	info, err := xcode_util.DefaultXcodeUtil.GetBuildTargetInfo(projectFile, projectType, scheme, arch)
	if err != nil {
		color.Error.Printf("获取构建目标信息失败: %v\n", err)
		return
	}

	color.Green.Println("🎯 构建目标信息:")
	fmt.Printf("   项目文件: %s\n", info.ProjectFile)
	fmt.Printf("   项目类型: %s\n", info.ProjectTypeName)
	fmt.Printf("   构建方案: %s\n", info.Scheme)
	fmt.Printf("   项目支持架构: %s\n", info.ProjectArchs)
	fmt.Printf("   构建目标架构: %s\n", info.TargetArch)
	fmt.Println()
}

// performBuild 执行构建
func performBuild(projectFile, projectType, scheme, buildPath, arch string, verbose, clean bool) error {
	color.Blue.Println("===========================================")
	color.Yellow.Println("🚀 开始构建过程...")
	color.Blue.Println("===========================================")
	fmt.Println()

	// 构建基础参数
	args := []string{}
	if projectType == "workspace" {
		args = append(args, "-workspace", projectFile)
	} else {
		args = append(args, "-project", projectFile)
	}

	args = append(args, "-scheme", scheme, "-configuration", "Release", "-derivedDataPath", buildPath)

	// 设置目标和架构
	args = append(args, "-destination", "generic/platform=macOS")
	if arch != "universal" {
		args = append(args, "ARCHS="+arch, "ONLY_ACTIVE_ARCH=NO")
	} else {
		args = append(args, "ARCHS=x86_64 arm64", "ONLY_ACTIVE_ARCH=NO")
	}

	// 添加静默参数
	if !verbose {
		args = append(args, "-quiet")
	}

	// 清理构建
	if clean {
		color.Yellow.Println("正在清理之前的构建...")
		cleanArgs := append(args, "clean")
		cleanCmd := exec.Command("xcodebuild", cleanArgs...)
		if verbose {
			cleanCmd.Stdout = os.Stdout
			cleanCmd.Stderr = os.Stderr
		}
		err := cleanCmd.Run()
		if err != nil {
			return fmt.Errorf("清理失败: %v", err)
		}
	}

	// 开始构建
	if arch == "universal" {
		color.Yellow.Println("开始构建应用 (通用二进制: x86_64 arm64)...")
	} else {
		color.Yellow.Printf("开始构建应用 (架构: %s)...\n", arch)
	}

	buildArgs := append(args, "build")
	buildCmd := exec.Command("xcodebuild", buildArgs...)

	if verbose {
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		fmt.Printf("执行命令: xcodebuild %s\n", strings.Join(buildArgs, " "))
		fmt.Println()
	}

	err := buildCmd.Run()
	if err != nil {
		return fmt.Errorf("构建失败: %v", err)
	}

	return nil
}
