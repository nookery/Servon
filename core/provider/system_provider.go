package provider

import (
	"bytes"
	"os"
	"os/exec"
	"servon/core/model"
	"servon/core/templates"
	"servon/core/utils/logger"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type serviceTemplateData struct {
	Command     string
	Args        string
	ServiceName string
}

// SystemProvider 系统管理器
type SystemProvider struct {
	RootCmd *cobra.Command
}

func NewSystemProvider() SystemProvider {
	return SystemProvider{
		RootCmd: &cobra.Command{},
	}
}

func (p *SystemProvider) InstallSoftware(name string, logChan chan<- string) error {
	// 检查操作系统类型
	osType := p.GetOSType()
	logger.InfoChan(logChan, "检测到操作系统: %s", osType)

	return nil
}

func (p *SystemProvider) UninstallSoftware(name string, logChan chan<- string) error {
	return nil
}

func (s *SystemProvider) GetOSType() OSType {
	// 尝试读取 /etc/os-release 文件
	cmd := exec.Command("cat", "/etc/os-release")
	output, err := cmd.Output()
	if err != nil {
		return model.Unknown
	}

	osInfo := strings.ToLower(string(output))

	switch {
	case strings.Contains(osInfo, "ubuntu"):
		return model.Ubuntu
	case strings.Contains(osInfo, "debian"):
		return model.Debian
	case strings.Contains(osInfo, "centos"):
		return model.CentOS
	case strings.Contains(osInfo, "redhat"):
		return model.RedHat
	default:
		return model.Unknown
	}
}

func (s *SystemProvider) CanUseApt() bool {
	osType := s.GetOSType()
	return osType == model.Ubuntu || osType == model.Debian
}

// RunBackgroundService 使用 systemd 在后台运行指定的命令作为服务
func (p *SystemProvider) RunBackgroundService(command string, args []string, logChan chan<- string) error {
	// 生成唯一的服务名称
	serviceName := "servon-" + strings.ReplaceAll(command, "/", "-") + ".service"

	// 准备模板数据
	data := serviceTemplateData{
		Command:     command,
		Args:        strings.Join(args, " "),
		ServiceName: serviceName,
	}

	// 读取并解析模板
	tmpl, err := template.ParseFS(templates.TemplateFS, "systemd_service.tmpl")
	if err != nil {
		logger.ErrorChan(logChan, "解析服务模板失败: %v", err)
		return err
	}

	// 渲染模板到缓冲区
	var serviceContent bytes.Buffer
	if err := tmpl.Execute(&serviceContent, data); err != nil {
		logger.ErrorChan(logChan, "渲染服务模板失败: %v", err)
		return err
	}

	// 确保日志目录存在
	if err := exec.Command("mkdir", "-p", "/var/log/servon").Run(); err != nil {
		logger.ErrorChan(logChan, "创建日志目录失败: %v", err)
		return err
	}

	// 写入服务文件
	serviceFilePath := "/etc/systemd/system/" + serviceName
	if err := os.WriteFile(serviceFilePath, serviceContent.Bytes(), 0644); err != nil {
		logger.ErrorChan(logChan, "创建服务文件失败: %v", err)
		return err
	}

	logger.InfoChan(logChan, "正在启动服务: %s", serviceName)

	// 重新加载 systemd 配置
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		logger.ErrorChan(logChan, "重载 systemd 配置失败: %v", err)
		return err
	}

	// 启动并启用服务
	if err := exec.Command("systemctl", "enable", "--now", serviceName).Run(); err != nil {
		logger.ErrorChan(logChan, "启动服务失败: %v", err)
		return err
	}

	logger.InfoChan(logChan, "服务已成功启动: %s", serviceName)
	logger.InfoChan(logChan, "可以通过以下命令查看日志:")
	logger.InfoChan(logChan, "  journalctl -u %s -f", serviceName)
	logger.InfoChan(logChan, "或查看日志文件:")
	logger.InfoChan(logChan, "  tail -f /var/log/servon/%s.log", serviceName)

	return nil
}

// StopBackgroundService 停止后台运行的服务
func (p *SystemProvider) StopBackgroundService(command string, logChan chan<- string) error {
	serviceName := "servon-" + strings.ReplaceAll(command, "/", "-") + ".service"

	logger.InfoChan(logChan, "正在停止服务: %s", serviceName)

	// 停止并禁用服务
	if err := exec.Command("systemctl", "disable", "--now", serviceName).Run(); err != nil {
		logger.ErrorChan(logChan, "停止服务失败: %v", err)
		return err
	}

	// 删除服务文件
	if err := os.Remove("/etc/systemd/system/" + serviceName); err != nil {
		logger.ErrorChan(logChan, "删除服务文件失败: %v", err)
		return err
	}

	// 重新加载 systemd 配置
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		logger.ErrorChan(logChan, "重载 systemd 配置失败: %v", err)
		return err
	}

	logger.InfoChan(logChan, "服务已成功停止并移除: %s", serviceName)
	return nil
}
