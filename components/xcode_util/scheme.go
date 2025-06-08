package xcode_util

import (
	"os/exec"
	"strings"
)

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
