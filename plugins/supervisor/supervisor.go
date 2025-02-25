package supervisor

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"servon/core"
	"strings"
)

type SupervisorPlugin struct {
	info core.SoftwareInfo
	*core.App
}

func Setup(app *core.App) {
	supervisor := NewSupervisorPlugin(app)
	app.RegisterSoftware("supervisor", supervisor)
}

func NewSupervisorPlugin(app *core.App) core.SuperSoft {
	return &SupervisorPlugin{
		App: app,
		info: core.SoftwareInfo{
			Name:        "supervisor",
			Description: "Supervisor 是运行时环境",
		},
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
