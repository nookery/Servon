package xcode

import (
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"

	"servon/components/dmg_util"
	"servon/components/xcode_util"
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "创建 DMG 安装包",
	Long:  color.Success.Render("\r\n为 macOS 应用程序创建 DMG 安装包"),
	Run: func(cmd *cobra.Command, args []string) {
		appPath, _ := cmd.Flags().GetString("path")
		outputDir, _ := cmd.Flags().GetString("output")
		verbose, _ := cmd.Flags().GetBool("verbose")

		// 检查应用程序路径是否指定
		if appPath == "" {
			color.Warnln("❌ 错误: 未指定应用程序路径")
			color.Info.Println("💡 使用方法: go run main.go xcode package --path /path/to/YourApp.app")
			os.Exit(1)
		}

		// 检查应用程序是否存在
		if _, err := os.Stat(appPath); os.IsNotExist(err) {
			color.Error.Printf("❌ 应用程序不存在: %s\n", appPath)
			os.Exit(1)
		}

		// 检查是否为 .app 文件
		if !strings.HasSuffix(appPath, ".app") {
			color.Error.Println("❌ 错误: 指定的路径不是 .app 应用程序")
			os.Exit(1)
		}

		// 设置默认输出目录
		if outputDir == "" {
			outputDir = "./temp"
		}

		// 显示配置信息
		showPackageConfig(appPath, outputDir, verbose)

		// 检查依赖
		err := dmg_util.DefaultDMGUtil.CheckDependencies()
		if err != nil {
			color.Error.Printf("❌ %s\n", err.Error())
			os.Exit(1)
		}

		// 显示应用信息
		xcode_util.DefaultXcodeUtil.ShowAppInfo(appPath)

		// 生成 DMG 文件名
		dmgFileName := generateDMGFileName(appPath)

		// 创建 DMG 文件
		dmgFile, err := dmg_util.DefaultDMGUtil.CreateDMG(appPath, outputDir, dmgFileName, verbose)
		if err != nil {
			color.Warnf("❌ 创建 DMG 失败: %s\n", err.Error())
			os.Exit(1)
		}

		// 显示结果
		dmg_util.DefaultDMGUtil.ShowResults(dmgFile)
	},
}

func init() {
	packageCmd.Flags().StringP("path", "p", "", "应用程序路径 (.app 文件)")
	packageCmd.Flags().StringP("output", "o", "./temp", "DMG 输出目录")
	packageCmd.Flags().BoolP("verbose", "v", false, "详细日志输出")
}

// showPackageConfig 显示配置信息
func showPackageConfig(appPath, outputDir string, verbose bool) {
	color.Blue.Println("===========================================")
	color.Blue.Println("         🚀 DMG 创建脚本                ")
	color.Blue.Println("===========================================")
	fmt.Println()

	color.Blue.Println("⚙️  配置信息")
	color.Info.Printf("应用程序: %s\n", appPath)
	color.Info.Printf("输出目录: %s\n", outputDir)
	color.Info.Printf("详细日志: %t\n", verbose)
	fmt.Println()
}

// generateDMGFileName 生成 DMG 文件名
// 格式: 应用名称-架构-版本
func generateDMGFileName(appPath string) string {
	// 获取应用名称
	appName := strings.TrimSuffix(strings.Split(appPath, "/")[len(strings.Split(appPath, "/"))-1], ".app")
	
	// 获取应用信息
	appInfo := xcode_util.DefaultXcodeUtil.GetAppInfo(appPath)
	
	// 清理版本号中的特殊字符
	version := strings.ReplaceAll(appInfo.Version, ".", "_")
	version = strings.ReplaceAll(version, " ", "_")
	
	// 清理架构名称中的特殊字符
	arch := strings.ReplaceAll(appInfo.Architecture, " ", "_")
	arch = strings.ReplaceAll(arch, ",", "+")
	
	// 生成文件名: 应用名称-架构-版本
	fileName := fmt.Sprintf("%s-%s-%s", appName, arch, version)
	
	return fileName
}
