package service

import (
	"os/exec"
	"servon/core/utils/logger"
	"strings"
)

func IsActive(serviceName string) bool {
	logger.Debug("检查服务状态: %s", serviceName)

	// 首先尝试 systemctl
	cmd := exec.Command("systemctl", "is-active", serviceName)
	output, err := cmd.CombinedOutput()
	if err == nil {
		return strings.TrimSpace(string(output)) == "active"
	}

	// 如果 systemctl 失败，尝试 service 命令
	cmd = exec.Command("service", serviceName, "status")
	output, err = cmd.CombinedOutput()
	if err != nil {
		logger.Debug("服务 %s 未运行: %v\n输出: %s", serviceName, err, string(output))
		return false
	}

	outputStr := string(output)
	return strings.Contains(outputStr, "is running") ||
		strings.Contains(outputStr, "start/running") ||
		strings.Contains(outputStr, "status=0/SUCCESS")
}
