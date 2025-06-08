package dmg_util

//
// DMG 工具组件
//
// 这个组件封装了与 macOS DMG 文件创建相关的工具函数，
// 包括 DMG 文件创建、挂载、压缩等功能。

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
)

var DefaultDMGUtil = &DMGUtil{}
var CreateDMG = DefaultDMGUtil.CreateDMG
var CheckDependencies = DefaultDMGUtil.CheckDependencies

type DMGUtil struct{}

// NewDMGUtil 创建新的 DMGUtil 实例
func NewDMGUtil() *DMGUtil {
	return &DMGUtil{}
}

// CreateDMG 创建 DMG 文件
// appPath: 应用程序路径
// outputDir: 输出目录
// dmgFileName: DMG 文件名
// verbose: 是否显示详细日志
// 返回: DMG 文件路径和错误信息
func (d *DMGUtil) CreateDMG(appPath, outputDir, dmgFileName string, verbose bool) (string, error) {
	color.Blue.Println("📦 创建 DMG 安装包")

	// 设置输出目录
	var dmgPath string
	if outputDir != "." {
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			return "", fmt.Errorf("无法创建输出目录: %v", err)
		}

		// 切换到输出目录
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(outputDir)
		
		// 如果是绝对路径，保持不变；如果是相对路径，添加 ../
		if !filepath.IsAbs(appPath) {
			appPath = "../" + appPath
		}
	}

	color.Info.Println("创建方法: hdiutil (原生)")
	relativeDMG, err := d.createDMGWithHdiutil(appPath, dmgFileName, verbose)
	if err != nil {
		return "", err
	}

	// 构建完整的 DMG 路径
	if outputDir != "." {
		dmgPath = filepath.Join(outputDir, relativeDMG)
	} else {
		dmgPath = relativeDMG
	}

	return dmgPath, nil
}

// CheckDependencies 检查依赖工具
func (d *DMGUtil) CheckDependencies() error {
	color.Blue.Println("🔍 检查依赖工具")

	// 检查 hdiutil
	if _, err := exec.LookPath("hdiutil"); err != nil {
		return fmt.Errorf("未找到 hdiutil 工具")
	}
	color.Success.Println("✅ hdiutil 可用")

	fmt.Println()
	return nil
}

// ShowResults 显示 DMG 创建结果
func (d *DMGUtil) ShowResults(dmgFile string) {
	color.Blue.Println("📋 DMG 创建结果")

	if _, err := os.Stat(dmgFile); err == nil {
		fileSize := "未知"
		if sizeOutput := d.getCommandOutput("ls", "-lh", dmgFile); sizeOutput != "" {
			parts := strings.Fields(sizeOutput)
			if len(parts) >= 5 {
				fileSize = parts[4]
			}
		}
		color.Info.Printf("%s: %s\n", dmgFile, fileSize)
	} else {
		color.Error.Printf("⚠️ 无法找到 DMG 文件: %s\n", dmgFile)
	}

	fmt.Println()
	color.Success.Println("✅ DMG 安装包创建完成！")
}

// createDMGWithHdiutil 使用 hdiutil 创建 DMG
func (d *DMGUtil) createDMGWithHdiutil(appPath, dmgFileName string, verbose bool) (string, error) {
	// 使用传入的文件名
	finalDMG := dmgFileName + ".dmg"
	tempDMG := dmgFileName + "-temp.dmg"

	// 清理已存在的文件
	os.Remove(tempDMG)
	os.Remove(finalDMG)

	// 创建临时 DMG
	// 从应用路径提取应用名称用于卷名
	appName := filepath.Base(appPath)
	appName = strings.TrimSuffix(appName, ".app")
	args := []string{"create", "-srcfolder", appPath, "-format", "UDRW", "-volname", appName, tempDMG}
	err := d.executeCommand("hdiutil", args, "创建临时 DMG", verbose)
	if err != nil {
		return "", err
	}

	// 挂载 DMG
	mountOutput, err := exec.Command("hdiutil", "attach", tempDMG, "-readwrite", "-noverify", "-noautoopen").Output()
	if err != nil {
		return "", fmt.Errorf("挂载 DMG 失败: %v", err)
	}

	// 解析挂载点
	mountPoint := ""
	lines := strings.Split(string(mountOutput), "\n")
	for _, line := range lines {
		if strings.Contains(line, "/Volumes/") {
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.HasPrefix(part, "/Volumes/") {
					mountPoint = part
					break
				}
			}
			break
		}
	}

	if mountPoint == "" {
		return "", fmt.Errorf("无法找到挂载点")
	}

	// 创建应用程序快捷方式
	err = d.executeCommand("ln", []string{"-s", "/Applications", filepath.Join(mountPoint, "Applications")}, "创建 Applications 快捷方式", verbose)
	if err != nil {
		// 卸载 DMG
		exec.Command("hdiutil", "detach", mountPoint).Run()
		return "", err
	}

	// 卸载 DMG
	err = d.executeCommand("hdiutil", []string{"detach", mountPoint}, "卸载 DMG", verbose)
	if err != nil {
		return "", err
	}

	// 压缩为最终文件名
	err = d.executeCommand("hdiutil", []string{"convert", tempDMG, "-format", "UDZO", "-imagekey", "zlib-level=9", "-o", finalDMG}, "压缩 DMG", verbose)
	if err != nil {
		return "", err
	}

	// 删除临时文件
	os.Remove(tempDMG)

	return finalDMG, nil
}

// executeCommand 执行命令
func (d *DMGUtil) executeCommand(command string, args []string, description string, verbose bool) error {
	if verbose {
		color.Blue.Printf("🔧 %s\n", description)
		color.Cyan.Printf("命令: %s %s\n", command, strings.Join(args, " "))
	}

	cmd := exec.Command(command, args...)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s失败: %v", description, err)
	}

	if verbose {
		color.Success.Printf("✅ %s 完成\n", description)
	}

	return nil
}

// getCommandOutput 获取命令输出
func (d *DMGUtil) getCommandOutput(command string, args ...string) string {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}