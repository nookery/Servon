package xcode

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var bumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "递增应用版本号",
	Long:  color.Success.Render("\r\n自动递增应用程序的修订版本号（最后一位数字）"),
	Run: func(cmd *cobra.Command, args []string) {
		projectFile, _ := cmd.Flags().GetString("project")
		workDir, _ := cmd.Flags().GetString("workdir")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// 显示环境信息
		showBumpEnvironmentInfo()

		// 如果指定了工作目录，切换到该目录
		if workDir != "" {
			if err := os.Chdir(workDir); err != nil {
				color.Error.Printf("❌ 无法切换到工作目录 %s: %s\n", workDir, err.Error())
				os.Exit(1)
			}
			color.Info.Printf("📁 切换到工作目录: %s\n", workDir)
		}

		// 如果没有指定项目文件，自动查找
		if projectFile == "" {
			cwd, _ := os.Getwd()
			color.Info.Printf("🔍 正在搜索目录: %s\n", cwd)
			var err error
			projectFile, err = findPbxprojFile()
			if err != nil {
				color.Warnf("❌ %s\n", err.Error())
				color.Info.Println("💡 使用方法: servon xcode bump -p /path/to/project.pbxproj 或 -w /path/to/workdir")
				os.Exit(0)
			}
			color.Success.Printf("✅ 找到项目文件: %s\n", projectFile)
		}

		color.Info.Printf("📁 项目文件: %s\n", projectFile)

		// 获取当前版本号
		currentVersion, err := getVersionFromProject(projectFile)
		if err != nil {
			color.Error.Printf("❌ %s\n", err.Error())
			os.Exit(2)
		}

		color.Info.Printf("📱 当前版本: %s\n", currentVersion)

		// 计算新版本号
		newVersion, err := incrementVersion(currentVersion)
		if err != nil {
			color.Error.Printf("❌ %s\n", err.Error())
			os.Exit(3)
		}

		color.Success.Printf("🚀 新版本: %s\n", newVersion)

		if dryRun {
			color.Yellow.Println("🔍 预览模式，不会实际修改文件")
			return
		}

		// 更新项目文件
		err = updateVersionInProject(projectFile, currentVersion, newVersion)
		if err != nil {
			color.Error.Printf("❌ 更新版本失败: %s\n", err.Error())
			os.Exit(4)
		}

		color.Success.Println("✅ 版本号更新成功！")

		// 显示 Git 状态
		showGitStatus()
	},
}

func init() {
	bumpCmd.Flags().StringP("project", "p", "", "指定 .pbxproj 文件路径")
	bumpCmd.Flags().StringP("workdir", "w", "", "指定工作目录，在该目录中搜索项目文件")
	bumpCmd.Flags().Bool("dry-run", false, "预览模式，不实际修改文件")
}

// incrementVersion 递增版本号的最后一位
func incrementVersion(version string) (string, error) {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("版本号格式不正确，期望格式: x.y.z")
	}

	// 解析最后一位数字
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", fmt.Errorf("无法解析版本号的修订版本: %v", err)
	}

	// 递增
	patch++
	parts[2] = strconv.Itoa(patch)

	return strings.Join(parts, "."), nil
}

// updateVersionInProject 更新项目文件中的版本号
func updateVersionInProject(projectFile, oldVersion, newVersion string) error {
	content, err := os.ReadFile(projectFile)
	if err != nil {
		return fmt.Errorf("无法读取项目文件: %v", err)
	}

	// 替换版本号
	oldPattern := fmt.Sprintf("MARKETING_VERSION = %s", oldVersion)
	newPattern := fmt.Sprintf("MARKETING_VERSION = %s", newVersion)
	newContent := strings.ReplaceAll(string(content), oldPattern, newPattern)

	// 写回文件
	err = os.WriteFile(projectFile, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("无法写入项目文件: %v", err)
	}

	return nil
}
