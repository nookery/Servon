package service

import "servon/cmd/internal/utils"
import "os/exec"
import "fmt"

func Reload(serviceName string) error {
	utils.Info("正在重载服务(container): %s", serviceName)

	cmd := exec.Command("service", serviceName, "reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.Error("重载服务失败 %s: %v\n输出: %s", serviceName, err, string(output))
		return fmt.Errorf("重载服务失败: %v", err)
	}

	utils.Info("服务已成功重载: %s", serviceName)
	return nil
}
