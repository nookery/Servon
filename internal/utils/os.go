package utils

import (
	"os/exec"
	"strings"
)

// OSType 表示操作系统类型
type OSType string

const (
	Ubuntu  OSType = "ubuntu"
	Debian  OSType = "debian"
	CentOS  OSType = "centos"
	RedHat  OSType = "redhat"
	Unknown OSType = "unknown"
)

// GetOSType 获取当前操作系统类型
func GetOSType() OSType {
	// 尝试读取 /etc/os-release 文件
	cmd := exec.Command("cat", "/etc/os-release")
	output, err := cmd.Output()
	if err != nil {
		return Unknown
	}

	osInfo := strings.ToLower(string(output))

	switch {
	case strings.Contains(osInfo, "ubuntu"):
		return Ubuntu
	case strings.Contains(osInfo, "debian"):
		return Debian
	case strings.Contains(osInfo, "centos"):
		return CentOS
	case strings.Contains(osInfo, "redhat"):
		return RedHat
	default:
		return Unknown
	}
}
