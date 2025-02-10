package libs

import (
	"os"
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

// GetOSType 获取操作系统类型
func GetOSType() OSType {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return Unknown
	}

	content := strings.ToLower(string(data))

	switch {
	case strings.Contains(content, "ubuntu"):
		return Ubuntu
	case strings.Contains(content, "debian"):
		return Debian
	case strings.Contains(content, "centos"):
		return CentOS
	case strings.Contains(content, "red hat"):
		return RedHat
	default:
		return Unknown
	}
}
