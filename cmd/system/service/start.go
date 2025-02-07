package service

import (
	"fmt"
	"os/exec"
	"servon/cmd/utils"
)

func Start(serviceName string) error {
	utils.Debug("正在启动服务: %s", serviceName)

	// 尝试使用 systemctl
	cmd := exec.Command("systemctl", "start", serviceName)
	err := utils.StreamCommand(cmd)
	if err != nil {
		// 如果 systemctl 失败，尝试使用 service 命令
		cmd = exec.Command("service", serviceName, "start")
		err = utils.StreamCommand(cmd)
		if err != nil {
			return fmt.Errorf("启动服务失败: %v", err)
		}
	}

	// 验证服务是否成功启动
	if !IsActive(serviceName) {
		errMsg := fmt.Sprintf("%s (服务未能成功运行)", serviceName)
		utils.Error("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	utils.Info("服务已成功启动: %s", serviceName)
	return nil
}
