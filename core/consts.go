package core

import (
	"os"
	"path/filepath"
	"servon/components/dev_util"
)

var DataRootFolder = getDataRootFolder()

const LoggerFolder = "/logs"
const DefaultHost = "0.0.0.0"
const DefaultPort = 8080

func getDataRootFolder() string {
	// 使用 dev_util 组件检测是否为开发模式
	if dev_util.DefaultDevUtil.IsDev() {
		// 开发模式：使用当前工作目录
		wd, err := os.Getwd()
		if err != nil {
			return "/data"
		}
		return wd
	} else {
		// 生产模式：使用用户家目录
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// 如果获取用户家目录失败，回退到默认路径
			return "/data"
		}
		return filepath.Join(homeDir, ".servon")
	}
}

const (
	Ubuntu  OSType = "ubuntu"
	Debian  OSType = "debian"
	CentOS  OSType = "centos"
	RedHat  OSType = "redhat"
	Unknown OSType = "unknown"
)
