package utils

import (
	"os/exec"
	"servon/core/model"
	"strings"
)

type OSType = model.OSType

// GetOSType 获取当前操作系统类型
func GetOSType() OSType {
	// 尝试读取 /etc/os-release 文件
	cmd := exec.Command("cat", "/etc/os-release")
	output, err := cmd.Output()
	if err != nil {
		return model.Unknown
	}

	osInfo := strings.ToLower(string(output))

	switch {
	case strings.Contains(osInfo, "ubuntu"):
		return model.Ubuntu
	case strings.Contains(osInfo, "debian"):
		return model.Debian
	case strings.Contains(osInfo, "centos"):
		return model.CentOS
	case strings.Contains(osInfo, "redhat"):
		return model.RedHat
	default:
		return model.Unknown
	}
}
