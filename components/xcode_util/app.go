package xcode_util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
)

// GetAppInfo 获取应用信息
func (x *XcodeUtil) GetAppInfo(appPath string) AppInfo {
	info := AppInfo{
		Version:      "unknown",
		Architecture: "unknown",
	}

	// 获取应用版本
	infoPath := appPath + "/Contents/Info.plist"
	if _, err := os.Stat(infoPath); err == nil {
		if version := x.GetCommandOutput("plutil", "-p", infoPath); version != "" {
			lines := strings.Split(version, "\n")
			for _, line := range lines {
				if strings.Contains(line, "CFBundleShortVersionString") {
					parts := strings.Split(line, `"`)
					if len(parts) >= 4 {
						info.Version = parts[3]
						break
					}
				}
			}
		}
	}

	// 获取应用架构
	appName := strings.TrimSuffix(strings.Split(appPath, "/")[len(strings.Split(appPath, "/"))-1], ".app")
	executablePath := appPath + "/Contents/MacOS/" + appName

	if _, err := os.Stat(executablePath); err == nil {
		archInfo := x.GetCommandOutput("lipo", "-info", executablePath)
		if archInfo != "" {
			if strings.Contains(archInfo, "Non-fat file") {
				// 单一架构
				if strings.Contains(archInfo, "architecture: ") {
					parts := strings.Split(archInfo, "architecture: ")
					if len(parts) >= 2 {
						info.Architecture = strings.TrimSpace(parts[1])
					}
				}
			} else if strings.Contains(archInfo, "Architectures in the fat file") {
				// 通用二进制文件 (多架构)
				lines := strings.Split(archInfo, "\n")
				for _, line := range lines {
					if strings.Contains(line, "are: ") {
						parts := strings.Split(line, "are: ")
						if len(parts) >= 2 {
							archs := strings.TrimSpace(parts[1])
							// 将多架构格式化为简洁形式
							if strings.Contains(archs, "x86_64") && strings.Contains(archs, "arm64") {
								info.Architecture = "universal"
							} else {
								info.Architecture = strings.ReplaceAll(archs, " ", "+")
							}
						}
						break
					}
				}
			}
		}
	}

	return info
}

// ShowAppInfo 显示应用程序信息
func (x *XcodeUtil) ShowAppInfo(appPath string) {
	color.Greenln("🎯 应用程序信息:")
	fmt.Printf("   应用路径: %s\n", appPath)

	// 读取 Info.plist
	infoPath := filepath.Join(appPath, "Contents/Info.plist")
	if _, err := os.Stat(infoPath); err == nil {
		if version := x.GetCommandOutput("plutil", "-p", infoPath); version != "" {
			lines := strings.Split(version, "\n")
			for _, line := range lines {
				if strings.Contains(line, "CFBundleShortVersionString") {
					parts := strings.Split(line, `"`)
					if len(parts) >= 4 {
						fmt.Printf("   应用版本: %s\n", parts[3])
					}
				} else if strings.Contains(line, "CFBundleVersion") {
					parts := strings.Split(line, `"`)
					if len(parts) >= 4 {
						fmt.Printf("   构建版本: %s\n", parts[3])
					}
				} else if strings.Contains(line, "CFBundleIdentifier") {
					parts := strings.Split(line, `"`)
					if len(parts) >= 4 {
						fmt.Printf("   Bundle ID: %s\n", parts[3])
					}
				}
			}
		}
		// 从应用路径中提取应用名称
		appName := filepath.Base(appPath)
		if strings.HasSuffix(appName, ".app") {
			appName = strings.TrimSuffix(appName, ".app")
		}
		fmt.Printf("   应用名称: %s\n", appName)

		// 获取应用架构信息
		x.showAppArchitecture(appPath)
	} else {
		color.Warnln("   ⚠️  无法读取应用信息")
	}
	fmt.Println()
}

// showAppArchitecture 显示应用架构信息
func (x *XcodeUtil) showAppArchitecture(appPath string) {
	// 构建可执行文件路径
	appName := filepath.Base(appPath)
	if strings.HasSuffix(appName, ".app") {
		appName = strings.TrimSuffix(appName, ".app")
	}
	executablePath := filepath.Join(appPath, "Contents/MacOS", appName)

	// 检查可执行文件是否存在
	if _, err := os.Stat(executablePath); err != nil {
		fmt.Printf("   应用架构: 未知 (无法找到可执行文件)\n")
		return
	}

	// 使用 lipo 命令获取架构信息
	archInfo := x.GetCommandOutput("lipo", "-info", executablePath)
	if archInfo == "" {
		fmt.Printf("   应用架构: 未知 (无法获取架构信息)\n")
		return
	}

	// 解析架构信息
	var architectures string
	if strings.Contains(archInfo, "Non-fat file") {
		// 单一架构
		if strings.Contains(archInfo, "architecture: ") {
			parts := strings.Split(archInfo, "architecture: ")
			if len(parts) >= 2 {
				architectures = strings.TrimSpace(parts[1])
			}
		}
	} else if strings.Contains(archInfo, "Architectures in the fat file") {
		// 通用二进制文件 (多架构)
		lines := strings.Split(archInfo, "\n")
		for _, line := range lines {
			if strings.Contains(line, "are: ") {
				parts := strings.Split(line, "are: ")
				if len(parts) >= 2 {
					architectures = strings.TrimSpace(parts[1])
					architectures = strings.ReplaceAll(architectures, " ", ", ")
				}
				break
			}
		}
	}

	if architectures != "" {
		fmt.Printf("   应用架构: %s\n", architectures)
	} else {
		fmt.Printf("   应用架构: 未知\n")
	}
}
