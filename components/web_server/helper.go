package web_server

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// findProcessByPort 通过端口查找进程ID
func findProcessByPort(port int) (int32, error) {
	return findProcessByPortWithVerbose(port, false, nil)
}

// findProcessByPortWithVerbose 通过端口查找进程ID（支持详细日志）
func findProcessByPortWithVerbose(port int, verbose bool, logger Logger) (int32, error) {
	if verbose && logger != nil {
		logger.Infof("🔍 开始检查端口 %d 是否被占用...", port)
	}

	// 使用更可靠的命令组合
	cmdStr := fmt.Sprintf("lsof -ti :%d | head -1", port)
	if verbose && logger != nil {
		logger.Infof("📋 执行命令: %s", cmdStr)
	}

	cmd := exec.Command("sh", "-c", cmdStr)
	output, err := cmd.Output()
	if err != nil {
		if verbose && logger != nil {
			logger.Infof("❌ 命令执行失败: %v", err)
			logger.Infof("✅ 端口 %d 未被占用，可以使用", port)
		}
		return 0, fmt.Errorf("查找进程失败: %v", err)
	}

	outputStr := strings.TrimSpace(string(output))
	if verbose && logger != nil {
		logger.Infof("📤 命令输出: '%s'", outputStr)
	}

	if outputStr == "" {
		if verbose && logger != nil {
			logger.Infof("✅ 端口 %d 未被占用，可以使用", port)
		}
		return 0, fmt.Errorf("端口 %d 上没有运行的进程", port)
	}

	pid, err := strconv.Atoi(outputStr)
	if err != nil {
		if verbose && logger != nil {
			logger.Infof("❌ 解析PID失败: %v", err)
		}
		return 0, fmt.Errorf("解析PID失败: %v", err)
	}

	if verbose && logger != nil {
		logger.Infof("🔍 发现进程 PID: %d，正在获取进程详细信息...", pid)
		// 获取进程详细信息
		getProcessDetails(int32(pid), verbose, logger)
	}

	return int32(pid), nil
}

// getProcessDetails 获取进程详细信息
func getProcessDetails(pid int32, verbose bool, logger Logger) {
	if !verbose || logger == nil {
		return
	}

	// 获取进程名称和命令行
	cmdStr := fmt.Sprintf("ps -p %d -o pid,ppid,comm,args -h", pid)
	logger.Infof("📋 获取进程详情命令: %s", cmdStr)

	cmd := exec.Command("sh", "-c", cmdStr)
	output, err := cmd.Output()
	if err != nil {
		logger.Infof("❌ 获取进程详情失败: %v", err)
		return
	}

	logger.Infof("📊 进程详情: \n%s", output)
}
