package libs

import (
	"fmt"
	"os/exec"
	"strings"
)

func IsActive(serviceName string) bool {
	Debug("检查服务状态: %s", serviceName)

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
		Debug("服务 %s 未运行: %v\n输出: %s", serviceName, err, string(output))
		return false
	}

	outputStr := string(output)
	return strings.Contains(outputStr, "is running") ||
		strings.Contains(outputStr, "start/running") ||
		strings.Contains(outputStr, "status=0/SUCCESS")
}

func Reload(serviceName string) error {
	Info("正在重载服务(container): %s", serviceName)

	cmd := exec.Command("service", serviceName, "reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		Error("重载服务失败 %s: %v\n输出: %s", serviceName, err, string(output))
		return fmt.Errorf("重载服务失败: %v", err)
	}

	Info("服务已成功重载: %s", serviceName)
	return nil

}

func Start(serviceName string) error {
	Debug("正在启动服务: %s", serviceName)

	// 尝试使用 systemctl
	cmd := exec.Command("systemctl", "start", serviceName)
	err := StreamCommand(cmd)
	if err != nil {
		// 如果 systemctl 失败，尝试使用 service 命令
		cmd = exec.Command("service", serviceName, "start")
		err = StreamCommand(cmd)
		if err != nil {
			return fmt.Errorf("启动服务失败: %v", err)
		}
	}

	// 验证服务是否成功启动
	if !IsActive(serviceName) {
		errMsg := fmt.Sprintf("%s (服务未能成功运行)", serviceName)
		Error("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	Info("服务已成功启动: %s", serviceName)
	return nil
}

func Stop(serviceName string) error {
	Info("正在停止服务(container): %s", serviceName)
	cmd := exec.Command("service", serviceName, "stop")
	output, err := cmd.CombinedOutput()
	if err != nil {
		Error("停止服务失败 %s: %v\n输出: %s", serviceName, err, string(output))
		return fmt.Errorf("停止服务失败: %v", err)
	}

	// 验证服务是否已停止
	if IsActive(serviceName) {
		errMsg := fmt.Sprintf("服务停止失败: %s (服务仍在运行)", serviceName)
		Error(errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	Info("服务已成功停止: %s", serviceName)
	return nil
}

// RunBackgroundService 使用 systemd 在后台运行指定的命令作为服务
func RunBackgroundService(command string, args []string, logChan chan<- string) (string, error) {
	return "", nil
}

// StopBackgroundService 停止并移除后台运行的服务
func StopBackgroundService(command string, logChan chan<- string) error {
	return nil
}
