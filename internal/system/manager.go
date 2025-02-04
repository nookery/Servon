package system

import (
	"os/exec"
	"strings"
)

// NewServiceManager returns appropriate service manager based on the environment
func NewServiceManager() ServiceManager {
	// 检查是否可以使用 systemctl
	cmd := exec.Command("systemctl", "--version")
	output, err := cmd.CombinedOutput()

	// 如果命令执行成功且输出中不包含容器相关的错误信息
	if err == nil && !strings.Contains(string(output), "not running in this container") {
		return NewSystemCtl()
	}

	// 如果 systemctl 不可用或在容器中，使用 service 命令
	return NewServiceCmd()
}
