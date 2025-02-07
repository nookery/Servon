package system

import (
	"os"
	"os/exec"
	"strings"
)

// GetOSInfo 获取操作系统信息
func GetOSInfo() (string, error) {
	// 尝试读取 /etc/os-release
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		lines := strings.Split(string(data), "\n")
		var name, version string
		for _, line := range lines {
			if strings.HasPrefix(line, "NAME=") {
				name = strings.Trim(strings.TrimPrefix(line, "NAME="), "\"")
			} else if strings.HasPrefix(line, "VERSION_ID=") {
				version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
			}
		}
		if name != "" {
			osInfo := name
			if version != "" {
				osInfo += " " + version
			}
			return osInfo, nil
		}
	}

	// 如果上面的方法失败，尝试使用 lsb_release 命令
	if out, err := exec.Command("lsb_release", "-ds").Output(); err == nil {
		return strings.TrimSpace(string(out)), nil
	}

	// 如果还是没有获取到，使用 uname
	if out, err := exec.Command("uname", "-sr").Output(); err == nil {
		return strings.TrimSpace(string(out)), nil
	}

	return "Unknown", nil
}
