package supervisor

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"servon/core"
	"strings"
	"time"
)

type SupervisorPlugin struct {
	info core.SoftwareInfo
	*core.App
	configDir string
}

func Setup(app *core.App) {
	supervisor := NewSupervisorPlugin(app)
	app.RegisterSoftware("supervisor", supervisor)
	app.RegisterService("supervisor", supervisor)
}

func NewSupervisorPlugin(app *core.App) core.SuperService {
	return &SupervisorPlugin{
		App: app,
		info: core.SoftwareInfo{
			Name:        "supervisor",
			Description: "Supervisor 是运行时环境",
		},
		configDir: "/etc/supervisor/conf.d",
	}
}

func (s *SupervisorPlugin) Install() error {
	osType := s.GetOSType()

	switch osType {
	case core.Ubuntu, core.Debian:
		s.SoftwareLogger.Infof("使用 APT 包管理器安装...")

		// 安装 supervisor
		if err := s.AptInstall("supervisor"); err != nil {
			return err
		}

	case core.CentOS, core.RedHat:
		errMsg := "暂不支持在 RHEL 系统上安装 Supervisor"
		return s.SoftwareLogger.LogAndReturnErrorf("%s", errMsg)

	default:
		errMsg := fmt.Sprintf("不支持的操作系统: %s", osType)
		return s.SoftwareLogger.LogAndReturnErrorf("%s", errMsg)
	}

	// 验证安装结果
	if !s.IsInstalled("supervisor") {
		errMsg := "Supervisor 安装验证失败，未检测到已安装的包"
		return s.SoftwareLogger.LogAndReturnErrorf("%s", errMsg)
	}

	s.SoftwareLogger.Success("Supervisor 安装完成")

	s.Start()

	return nil
}

// Uninstall 卸载
func (s *SupervisorPlugin) Uninstall() error {
	// 卸载软件包及其依赖
	if err := s.AptRemove("supervisor"); err != nil {
		return err
	}

	// 清理配置文件
	if err := s.AptPurge("supervisor"); err != nil {
		return err
	}

	// 清理自动安装的依赖
	err, _ := s.RunShell("apt-get", "autoremove", "-y")
	if err != nil {
		return fmt.Errorf("清理依赖失败:\n%s", err)
	}

	s.SoftwareLogger.Success("Supervisor 卸载完成")
	return nil
}

const (
	StatusNotInstalled = "not_installed" // 未安装
	StatusInstalled    = "installed"     // 已安装
	StatusRunning      = "running"       // 运行中
	StatusStopped      = "stopped"       // 已停止
	StatusError        = "error"         // 错误状态
)

func (s *SupervisorPlugin) GetStatus() (map[string]string, error) {
	s.SoftwareLogger.Infof("获取 Supervisor 状态")

	// 1. 检查是否安装
	if !s.IsInstalled("supervisor") {
		return map[string]string{
			"status":  StatusNotInstalled,
			"version": "",
			"message": "Supervisor 未安装",
		}, nil
	}

	// 2. 获取版本信息
	version := ""
	verCmd := exec.Command("supervisord", "--version")
	if verOutput, err := verCmd.CombinedOutput(); err == nil {
		version = strings.TrimSpace(string(verOutput))
	}

	// 3. 检查进程是否运行
	cmd := exec.Command("pgrep", "supervisord")
	if err := cmd.Run(); err != nil {
		return map[string]string{
			"status":  StatusStopped,
			"version": version,
			"message": "Supervisor 已安装但未运行",
		}, nil
	}

	// 4. 检查 socket 文件
	if _, err := os.Stat("/var/run/supervisor.sock"); os.IsNotExist(err) {
		return map[string]string{
			"status":  StatusError,
			"version": version,
			"message": "Supervisor 进程运行中但 socket 文件不存在",
		}, nil
	}

	// 5. 检查服务状态
	statusCmd := exec.Command("supervisorctl", "status")
	output, err := statusCmd.CombinedOutput()
	if err != nil {
		return map[string]string{
			"status":  StatusError,
			"version": version,
			"message": fmt.Sprintf("Supervisor 运行异常: %v", err),
			"error":   string(output),
		}, nil
	}

	// 6. 获取运行中的服务数量
	runningServices := 0
	totalServices := 0
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		totalServices++
		if strings.Contains(line, "RUNNING") {
			runningServices++
		}
	}

	// 7. 检查配置文件
	configExists := false
	configPaths := []string{
		"/etc/supervisor/supervisord.conf",
		"/etc/supervisord.conf",
	}
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configExists = true
			break
		}
	}

	return map[string]string{
		"status":           StatusRunning,
		"version":          version,
		"message":          "Supervisor 运行正常",
		"config_exists":    fmt.Sprintf("%v", configExists),
		"total_services":   fmt.Sprintf("%d", totalServices),
		"running_services": fmt.Sprintf("%d", runningServices),
		"pid_file":         s.getPidFile(),
		"uptime":           s.getUptime(),
	}, nil
}

// 获取 PID 文件路径和状态
func (s *SupervisorPlugin) getPidFile() string {
	pidPaths := []string{
		"/var/run/supervisord.pid",
		"/tmp/supervisord.pid",
	}

	for _, path := range pidPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return "未找到"
}

// 获取运行时间
func (s *SupervisorPlugin) getUptime() string {
	cmd := exec.Command("ps", "-p", s.getSupervisorPid(), "-o", "etime=")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "未知"
	}
	return strings.TrimSpace(string(output))
}

// 获取 Supervisor 进程 PID
func (s *SupervisorPlugin) getSupervisorPid() string {
	cmd := exec.Command("pgrep", "supervisord")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// 检查配置文件语法
func (s *SupervisorPlugin) checkConfig() error {
	cmd := exec.Command("supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf", "-t")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("配置文件检查失败: %v\n%s", err, string(output))
	}
	return nil
}

func (s *SupervisorPlugin) GetInfo() core.SoftwareInfo {
	return s.info
}

func (s *SupervisorPlugin) Start() error {
	s.SoftwareLogger.Infof("Supervisor 开始启动")

	err, _ := s.RunShellWithSudo("supervisord", "-c", "/etc/supervisor/supervisord.conf")
	if err != nil {
		return err
	}

	s.SoftwareLogger.Success("Supervisor 启动成功")
	return nil
}

func (s *SupervisorPlugin) Stop() error {
	s.SoftwareLogger.Infof("Supervisor 开始停止")

	err, _ := s.RunShellWithSudo("supervisorctl", "shutdown")
	if err != nil {
		return err
	}

	s.SoftwareLogger.Success("Supervisor 停止成功")
	return nil
}

// 实现 Service 接口的方法

func (s *SupervisorPlugin) StartService(serviceName string) error {
	s.SoftwareLogger.Infof("启动服务: %s", serviceName)

	err, output := s.RunShellWithSudo("supervisorctl", "start", serviceName)
	if err != nil {
		return fmt.Errorf("启动服务 %s 失败: %v\n%s", serviceName, err, output)
	}

	s.SoftwareLogger.Successf("服务 %s 启动成功", serviceName)
	return nil
}

func (s *SupervisorPlugin) StopService(serviceName string) error {
	s.SoftwareLogger.Infof("停止服务: %s", serviceName)

	err, output := s.RunShellWithSudo("supervisorctl", "stop", serviceName)
	if err != nil {
		return fmt.Errorf("停止服务 %s 失败: %v\n%s", serviceName, err, output)
	}

	s.SoftwareLogger.Successf("服务 %s 停止成功", serviceName)
	return nil
}

func (s *SupervisorPlugin) AddBackgroundService(serviceName string, command string, args []string, env []string) (string, error) {
	s.SoftwareLogger.Infof("添加后台服务: %s", serviceName)

	// 确保配置目录存在
	if err := s.ensureConfigDir(); err != nil {
		return "", err
	}

	// 创建配置文件
	configPath := filepath.Join(s.configDir, serviceName+".conf")

	// 检查是否已存在
	if _, err := os.Stat(configPath); err == nil {
		return "", fmt.Errorf("服务配置文件已存在: %s", configPath)
	}

	// 构建完整命令
	fullCommand := command
	if len(args) > 0 {
		fullCommand += " " + strings.Join(args, " ")
	}

	// 构建配置内容
	configContent := fmt.Sprintf(`[program:%s]
command=%s
autostart=true
autorestart=true
stderr_logfile=/var/log/supervisor/%s.err.log
stdout_logfile=/var/log/supervisor/%s.out.log
`, serviceName, fullCommand, serviceName, serviceName)

	// 添加环境变量
	if len(env) > 0 {
		configContent += "environment="
		envPairs := make([]string, 0, len(env))
		for _, e := range env {
			envPairs = append(envPairs, e)
		}
		configContent += strings.Join(envPairs, ",") + "\n"
	}

	// 写入配置文件
	if err := ioutil.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return "", fmt.Errorf("写入配置文件失败: %v", err)
	}

	// 重新加载配置
	if err := s.reloadConfig(); err != nil {
		return "", err
	}

	// 启动服务
	if err := s.StartService(serviceName); err != nil {
		return "", err
	}

	return configPath, nil
}

func (s *SupervisorPlugin) StopBackgroundService(serviceName string) error {
	s.SoftwareLogger.Infof("停止后台服务: %s", serviceName)

	// 先停止服务
	if err := s.StopService(serviceName); err != nil {
		return err
	}

	// 删除配置文件
	configPath := filepath.Join(s.configDir, serviceName+".conf")
	if err := os.Remove(configPath); err != nil {
		return fmt.Errorf("删除配置文件失败: %v", err)
	}

	// 重新加载配置
	if err := s.reloadConfig(); err != nil {
		return err
	}

	s.SoftwareLogger.Successf("后台服务 %s 已停止并移除", serviceName)
	return nil
}

func (s *SupervisorPlugin) Restart() error {
	s.SoftwareLogger.Infof("重启 Supervisor")

	if err := s.Stop(); err != nil {
		return err
	}

	// 等待进程完全停止
	time.Sleep(2 * time.Second)

	return s.Start()
}

func (s *SupervisorPlugin) RestartService(serviceName string) error {
	s.SoftwareLogger.Infof("重启服务: %s", serviceName)

	err, output := s.RunShellWithSudo("supervisorctl", "restart", serviceName)
	if err != nil {
		return fmt.Errorf("重启服务 %s 失败: %v\n%s", serviceName, err, output)
	}

	s.SoftwareLogger.Successf("服务 %s 重启成功", serviceName)
	return nil
}

func (s *SupervisorPlugin) GetLogs(serviceName string, lines int) (string, error) {
	s.SoftwareLogger.Infof("获取服务 %s 的日志", serviceName)

	if lines <= 0 {
		lines = 100 // 默认获取100行
	}

	// 获取标准输出日志
	stdoutLogPath := fmt.Sprintf("/var/log/supervisor/%s.out.log", serviceName)
	stderrLogPath := fmt.Sprintf("/var/log/supervisor/%s.err.log", serviceName)

	// 检查日志文件是否存在
	if _, err := os.Stat(stdoutLogPath); os.IsNotExist(err) {
		return "", fmt.Errorf("服务 %s 的标准输出日志文件不存在", serviceName)
	}

	// 使用 tail 命令获取最后几行日志
	tailCmd := exec.Command("tail", "-n", fmt.Sprintf("%d", lines), stdoutLogPath)
	stdoutLog, err := tailCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取标准输出日志失败: %v", err)
	}

	// 获取错误日志
	var stderrLog []byte
	if _, err := os.Stat(stderrLogPath); err == nil {
		tailErrCmd := exec.Command("tail", "-n", fmt.Sprintf("%d", lines), stderrLogPath)
		stderrLog, _ = tailErrCmd.CombinedOutput()
	}

	// 组合日志
	result := fmt.Sprintf("=== 标准输出日志 (%s) ===\n%s\n\n", stdoutLogPath, string(stdoutLog))
	if len(stderrLog) > 0 {
		result += fmt.Sprintf("=== 错误日志 (%s) ===\n%s", stderrLogPath, string(stderrLog))
	}

	return result, nil
}

func (s *SupervisorPlugin) IsRunning(serviceName string) (bool, error) {
	s.SoftwareLogger.Infof("检查服务 %s 是否运行中", serviceName)

	err, output := s.RunShell("supervisorctl", "status", serviceName)
	if err != nil {
		if strings.Contains(output, "ERROR") && strings.Contains(output, "not found") {
			return false, fmt.Errorf("服务 %s 不存在", serviceName)
		}
		return false, fmt.Errorf("获取服务状态失败: %v", err)
	}

	return strings.Contains(output, "RUNNING"), nil
}

func (s *SupervisorPlugin) GetServiceConfig(serviceName string) (map[string]interface{}, error) {
	s.SoftwareLogger.Infof("获取服务 %s 的配置", serviceName)

	configPath := filepath.Join(s.configDir, serviceName+".conf")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("服务 %s 的配置文件不存在", serviceName)
	}

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析配置文件内容
	config := make(map[string]interface{})
	config["raw_config"] = string(content)

	// 提取关键配置项
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			config[key] = value
		}
	}

	return config, nil
}

func (s *SupervisorPlugin) UpdateServiceConfig(serviceName string, config map[string]interface{}) error {
	s.SoftwareLogger.Infof("更新服务 %s 的配置", serviceName)

	// 获取当前配置
	currentConfig, err := s.GetServiceConfig(serviceName)
	if err != nil {
		return err
	}

	// 如果提供了原始配置，直接使用
	if rawConfig, ok := config["raw_config"].(string); ok {
		configPath := filepath.Join(s.configDir, serviceName+".conf")
		if err := ioutil.WriteFile(configPath, []byte(rawConfig), 0644); err != nil {
			return fmt.Errorf("写入配置文件失败: %v", err)
		}
	} else {
		// 否则，更新各个配置项
		rawConfig := currentConfig["raw_config"].(string)
		lines := strings.Split(rawConfig, "\n")

		// 更新配置项
		for i, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine == "" || strings.HasPrefix(trimmedLine, ";") || strings.HasPrefix(trimmedLine, "#") {
				continue
			}

			parts := strings.SplitN(trimmedLine, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				if newValue, ok := config[key]; ok {
					lines[i] = fmt.Sprintf("%s=%v", key, newValue)
				}
			}
		}

		// 写入更新后的配置
		configPath := filepath.Join(s.configDir, serviceName+".conf")
		if err := ioutil.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
			return fmt.Errorf("写入配置文件失败: %v", err)
		}
	}

	// 重新加载配置
	if err := s.reloadConfig(); err != nil {
		return err
	}

	// 重启服务以应用新配置
	return s.RestartService(serviceName)
}

func (s *SupervisorPlugin) GetResourceUsage(serviceName string) (map[string]interface{}, error) {
	s.SoftwareLogger.Infof("获取服务 %s 的资源使用情况", serviceName)

	// 检查服务是否存在
	running, err := s.IsRunning(serviceName)
	if err != nil {
		return nil, err
	}

	if !running {
		return map[string]interface{}{
			"status": "stopped",
			"cpu":    0,
			"memory": 0,
		}, nil
	}

	// 获取服务的进程ID
	err, output := s.RunShell("supervisorctl", "pid", serviceName)
	if err != nil {
		return nil, fmt.Errorf("获取进程ID失败: %v", err)
	}

	pid := strings.TrimSpace(output)
	if pid == "" {
		return nil, fmt.Errorf("无法获取服务 %s 的进程ID", serviceName)
	}

	// 获取CPU和内存使用情况
	psCmd := exec.Command("ps", "-p", pid, "-o", "%cpu,%mem")
	psOutput, err := psCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("获取资源使用情况失败: %v", err)
	}

	// 解析输出
	lines := strings.Split(string(psOutput), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("无法解析资源使用情况输出")
	}

	// 解析CPU和内存使用率
	fields := strings.Fields(lines[1])
	if len(fields) < 2 {
		return nil, fmt.Errorf("无法解析CPU和内存使用率")
	}

	cpu := fields[0]
	memory := fields[1]

	return map[string]interface{}{
		"status": "running",
		"pid":    pid,
		"cpu":    cpu,
		"memory": memory,
	}, nil
}

func (s *SupervisorPlugin) GetServiceDetails(serviceName string) (map[string]interface{}, error) {
	s.SoftwareLogger.Infof("获取服务 %s 的详细信息", serviceName)

	// 获取服务状态
	err, output := s.RunShell("supervisorctl", "status", serviceName)
	if err != nil {
		if strings.Contains(output, "ERROR") && strings.Contains(output, "not found") {
			return nil, fmt.Errorf("服务 %s 不存在", serviceName)
		}
		return nil, fmt.Errorf("获取服务状态失败: %v", err)
	}

	// 解析状态输出
	details := make(map[string]interface{})
	details["raw_status"] = output

	// 解析状态行
	fields := strings.Fields(output)
	if len(fields) >= 2 {
		details["name"] = fields[0]
		details["status"] = fields[1]

		if len(fields) >= 4 && fields[1] == "RUNNING" {
			// 格式: name RUNNING pid uptime
			details["pid"] = fields[2]
			details["uptime"] = fields[3]
		}
	}

	// 获取配置信息
	config, err := s.GetServiceConfig(serviceName)
	if err == nil {
		details["config"] = config
	}

	// 获取资源使用情况
	resources, err := s.GetResourceUsage(serviceName)
	if err == nil {
		details["resources"] = resources
	}

	return details, nil
}

// 辅助方法

// 确保配置目录存在
func (s *SupervisorPlugin) ensureConfigDir() error {
	if _, err := os.Stat(s.configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(s.configDir, 0755); err != nil {
			return fmt.Errorf("创建配置目录失败: %v", err)
		}
	}
	return nil
}

// 重新加载配置
func (s *SupervisorPlugin) reloadConfig() error {
	err, output := s.RunShellWithSudo("supervisorctl", "reread")
	if err != nil {
		return fmt.Errorf("重新读取配置失败: %v\n%s", err, output)
	}

	err, output = s.RunShellWithSudo("supervisorctl", "update")
	if err != nil {
		return fmt.Errorf("更新配置失败: %v\n%s", err, output)
	}

	return nil
}
