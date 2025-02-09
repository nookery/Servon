package service

import (
	"fmt"
	"os/exec"
	"servon/utils/logger"
)

func Reload(serviceName string) error {
	logger.Info("正在重载服务(container): %s", serviceName)

	cmd := exec.Command("service", serviceName, "reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("重载服务失败 %s: %v\n输出: %s", serviceName, err, string(output))
		return fmt.Errorf("重载服务失败: %v", err)
	}

	logger.Info("服务已成功重载: %s", serviceName)
	return nil
}
