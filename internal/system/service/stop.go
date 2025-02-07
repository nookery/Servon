package service

import (
	"fmt"
	"os/exec"
	"servon/internal/utils"
)

func Stop(serviceName string) error {
	utils.Info("正在停止服务(container): %s", serviceName)
	cmd := exec.Command("service", serviceName, "stop")
	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.Error("停止服务失败 %s: %v\n输出: %s", serviceName, err, string(output))
		return fmt.Errorf("停止服务失败: %v", err)
	}

	// 验证服务是否已停止
	if IsActive(serviceName) {
		errMsg := fmt.Sprintf("服务停止失败: %s (服务仍在运行)", serviceName)
		utils.Error(errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	utils.Info("服务已成功停止: %s", serviceName)
	return nil
}
