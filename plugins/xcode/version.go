package xcode

import (
	"fmt"
	"os"
	"regexp"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "获取应用版本号",
	Long:  color.Success.Render("\r\n从 Xcode 项目配置文件中获取应用程序的营销版本号（MARKETING_VERSION）"),
	Run: func(cmd *cobra.Command, args []string) {
		projectFile, _ := cmd.Flags().GetString("project")
		workDir, _ := cmd.Flags().GetString("workdir")

		// 如果指定了工作目录，切换到该目录
		if workDir != "" {
			if err := os.Chdir(workDir); err != nil {
				color.Error.Printf("❌ 无法切换到工作目录 %s: %s\n", workDir, err.Error())
				os.Exit(1)
			}
		}

		// 如果没有指定项目文件，自动查找
		if projectFile == "" {
			cwd, _ := os.Getwd()
			color.Info.Printf("🔍 正在搜索目录: %s\n", cwd)
			color.Info.Println("💡 提示: 可使用 -p 参数指定 .pbxproj 文件路径，-w 参数指定工作目录")
			var err error
			projectFile, err = findPbxprojFile()
			if err != nil {
				color.Warnf("❌ %s\n", err.Error())
				color.Info.Println("💡 使用方法: servon xcode version -p /path/to/project.pbxproj 或 -w /path/to/workdir")
				os.Exit(0)
			}
			color.Success.Printf("✅ 找到项目文件: %s\n", projectFile)
		} else {
			color.Info.Printf("📁 使用指定的项目文件: %s\n", projectFile)
		}

		// 获取版本号
		version, err := getVersionFromProject(projectFile)
		if err != nil {
			color.Error.Printf("❌ %s\n", err.Error())
			os.Exit(2)
		}

		color.Success.Printf("📱 当前版本: %s\n", version)
	},
}

func init() {
	versionCmd.Flags().StringP("project", "p", "", "指定 .pbxproj 文件路径")
	versionCmd.Flags().StringP("workdir", "w", "", "指定工作目录，在该目录中搜索项目文件")
}

// getVersionFromProject 从项目文件中提取版本号
func getVersionFromProject(projectFile string) (string, error) {
	content, err := os.ReadFile(projectFile)
	if err != nil {
		return "", fmt.Errorf("无法读取项目文件: %v", err)
	}

	// 使用正则表达式查找 MARKETING_VERSION
	re := regexp.MustCompile(`MARKETING_VERSION\s*=\s*([0-9]+\.[0-9]+\.[0-9]+)`)
	matches := re.FindStringSubmatch(string(content))

	if len(matches) < 2 {
		return "", fmt.Errorf("未找到 MARKETING_VERSION")
	}

	return matches[1], nil
}
