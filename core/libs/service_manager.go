package libs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"servon/core/templates"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type ServiceManager struct {
	RootFolder string
	ConfigDir  string
}

type SupervisorConfig struct {
	ServiceName string
	Command     string
	Args        string
	RootFolder  string
	WorkingDir  string
	Environment string
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		RootFolder: DefaultDataManager.GetSoftwareRootFolder("supervisor"),
		ConfigDir:  DefaultDataManager.GetSoftwareRootFolder("supervisor") + "/conf.d",
	}
}

func (p *ServiceManager) checkSupervisorInstalled() error {
	cmd := exec.Command("which", "supervisord")
	if err := cmd.Run(); err != nil {
		printer.PrintErrorMessage("Supervisor未安装，请先安装Supervisor")
		PrintInfo("Ubuntu/Debian: sudo apt-get install supervisor")
		PrintInfo("CentOS/RHEL: sudo yum install supervisor")
		return fmt.Errorf("supervisor未安装")
	}
	return nil
}

func (p *ServiceManager) ensureConfigDir() error {
	if err := os.MkdirAll(p.ConfigDir, 0755); err != nil {
		return fmt.Errorf("创建Supervisor配置目录失败: %v", err)
	}

	// 确保日志目录存在
	logDir := filepath.Join(p.RootFolder, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建Supervisor日志目录失败: %v", err)
	}

	return nil
}

// getSupervisorConfigDir 获取supervisor配置目录
func (p *ServiceManager) getSupervisorConfigDir() (string, error) {
	// 常见的supervisor配置目录
	configDirs := []string{
		"/etc/supervisor/conf.d",
		"/etc/supervisord.d",
	}

	for _, dir := range configDirs {
		if _, err := os.Stat(dir); err == nil {
			return dir, nil
		}
	}

	return "", fmt.Errorf("未找到supervisor配置目录")
}

// createConfig 创建supervisor配置文件，返回配置文件路径
func (p *ServiceManager) createConfig(serviceName string, command string, args []string, envVars []string) (string, error) {
	if p.HasServiceConf(serviceName) {
		return "", fmt.Errorf("服务配置文件已存在: %s", serviceName)
	}

	configPath := filepath.Join(p.ConfigDir, serviceName+".conf")

	// 获取命令的绝对路径
	absCommand, err := exec.LookPath(command)
	if err != nil {
		return "", fmt.Errorf("找不到命令 %s: %v", command, err)
	}

	// 获取当前工作目录
	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取工作目录失败: %v", err)
	}

	// 获取模板内容
	tmplContent, err := templates.GetSupervisorConfigTemplate()
	if err != nil {
		return "", fmt.Errorf("获取模板内容失败: %v", err)
	}

	// 解析模板
	tmpl, err := template.New("supervisor").Parse(tmplContent)
	if err != nil {
		return "", fmt.Errorf("解析supervisor模板失败: %v", err)
	}

	// 准备模板数据
	config := SupervisorConfig{
		ServiceName: serviceName,
		Command:     absCommand,
		Args:        strings.Join(args, " "),
		RootFolder:  p.RootFolder,
		WorkingDir:  workingDir,
		Environment: strings.Join(envVars, ","),
	}

	// 创建配置文件
	f, err := os.Create(configPath)
	if err != nil {
		return "", fmt.Errorf("创建配置文件失败: %v", err)
	}
	defer f.Close()

	// 执行模板
	if err := tmpl.Execute(f, config); err != nil {
		return "", fmt.Errorf("生成配置文件失败: %v", err)
	}

	// 获取supervisor系统配置目录
	supervisorConfigDir, err := p.getSupervisorConfigDir()
	if err != nil {
		return "", fmt.Errorf("获取supervisor配置目录失败: %v", err)
	}

	// 创建软链接
	systemConfigPath := filepath.Join(supervisorConfigDir, serviceName+".conf")
	// 如果已存在，先删除
	if _, err := os.Lstat(systemConfigPath); err == nil {
		if err := os.Remove(systemConfigPath); err != nil {
			return "", fmt.Errorf("删除已存在的配置文件软链接失败: %v", err)
		}
	}

	// 创建软链接
	if err := os.Symlink(configPath, systemConfigPath); err != nil {
		return "", fmt.Errorf("创建配置文件软链接失败: %v", err)
	}

	PrintInfof("已创建配置文件软链接: %s -> %s", systemConfigPath, configPath)

	return configPath, nil
}

func (p *ServiceManager) checkSupervisorRunning() error {
	// 首先检查是否安装
	if err := p.checkSupervisorInstalled(); err != nil {
		return err
	}

	// 检查 supervisord 是否运行
	cmd := exec.Command("pgrep", "supervisord")
	if err := cmd.Run(); err != nil {
		PrintErrorMessage("Supervisor守护进程未运行")
		PrintCommandOutput("请使用以下命令启动 Supervisor:")
		PrintCommandOutput("supervisord -c /etc/supervisor/supervisord.conf")

		return fmt.Errorf("supervisor守护进程未运行")
	}

	// 检查 socket 文件
	if _, err := os.Stat("/var/run/supervisor.sock"); os.IsNotExist(err) {
		PrintErrorf("Supervisor socket 文件不存在")
		PrintInfo("请检查 Supervisor 配置文件是否正确")
		PrintInfo("配置文件通常位于: /etc/supervisor/supervisord.conf")
		return fmt.Errorf("supervisor socket 文件不存在")
	}

	return nil
}

// HasServiceConf 判断服务配置文件是否存在
func (p *ServiceManager) HasServiceConf(serviceName string) bool {
	if _, err := os.Stat(filepath.Join(p.ConfigDir, serviceName+".conf")); os.IsNotExist(err) {
		return false
	}

	return true
}

// GetServiceFilePath 获取服务配置文件路径
func (p *ServiceManager) GetServiceFilePath(serviceName string) string {
	return filepath.Join(p.ConfigDir, serviceName+".conf")
}

func (p *ServiceManager) IsActive(serviceName string) bool {
	if err := p.checkSupervisorInstalled(); err != nil {
		return false
	}

	cmd := exec.Command("supervisorctl", "status", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		PrintErrorf("检查服务状态失败 %s: %v", serviceName, err)
		return false
	}

	outputStr := string(output)
	return strings.Contains(outputStr, "RUNNING")
}

// Reload 重载服务
func (p *ServiceManager) Reload(serviceName string) error {
	if err := p.checkSupervisorInstalled(); err != nil {
		return err
	}

	PrintInfof("正在重载服务: %s", serviceName)

	// 首先执行 reread 命令读取新配置
	output, err := RunShellWithOutput("supervisorctl", "reread")
	if err != nil {
		return fmt.Errorf("读取配置失败: %v\n%s", err, output)
	}

	// 然后执行 update 命令更新配置
	output, err = RunShellWithOutput("supervisorctl", "update")
	if err != nil {
		return fmt.Errorf("更新配置失败: %v\n%s", err, output)
	}

	PrintSuccessf("服务已成功重载: %s", serviceName)
	return nil
}

// Start 启动服务
func (p *ServiceManager) Start(serviceName string) error {
	if err := p.checkSupervisorInstalled(); err != nil {
		return err
	}

	PrintInfof("正在启动服务: %s", serviceName)

	output, err := RunShellWithOutput("supervisorctl", "start", serviceName)
	if err != nil {
		return fmt.Errorf("启动服务失败: %v\n%s", err, output)
	}

	PrintSuccessf("服务已成功启动: %s", serviceName)
	return nil
}

// Stop 停止服务
func (p *ServiceManager) Stop(serviceName string) error {
	if err := p.checkSupervisorInstalled(); err != nil {
		return err
	}

	PrintInfof("正在停止服务: %s", serviceName)

	cmd := exec.Command("supervisorctl", "stop", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		PrintErrorf("停止服务失败 %s: %v\n输出: %s", serviceName, err, string(output))
		return fmt.Errorf("停止服务失败: %v", err)
	}

	// 验证服务是否已停止
	if p.IsActive(serviceName) {
		errMsg := fmt.Sprintf("服务停止失败: %s (服务仍在运行)", serviceName)
		PrintErrorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	PrintSuccessf("服务已成功停止: %s", serviceName)
	return nil
}

// AddBackgroundService 添加后台服务，返回配置文件路径
// serviceName: 服务名称
// command: 要执行的命令
// args: 命令参数
// env: 环境变量，格式如 ["KEY=VALUE", ...]
func (p *ServiceManager) AddBackgroundService(serviceName string, command string, args []string, env []string) (string, error) {
	printer.PrintInfof("正在添加后台服务: %s", command)

	if p.HasServiceConf(serviceName) {
		return "", fmt.Errorf("服务配置文件已存在: %s", serviceName)
	}

	if err := p.checkSupervisorInstalled(); err != nil {
		return "", err
	}

	if err := p.ensureConfigDir(); err != nil {
		return "", err
	}

	printer.PrintInfof("正在创建服务配置文件: %s", serviceName)

	configPath, err := p.createConfig(serviceName, command, args, env)
	if err != nil {
		return "", err
	}

	if err := p.Reload(serviceName); err != nil {
		return "", fmt.Errorf("重载supervisor配置失败: %v", err)
	}

	if err := p.Start(serviceName); err != nil {
		return "", err
	}

	return configPath, nil
}

// StopBackgroundService 停止后台服务
func (p *ServiceManager) StopBackgroundService(serviceName string, logChan chan<- string) error {
	if err := p.Stop(serviceName); err != nil {
		return err
	}

	configPath := filepath.Join(p.ConfigDir, serviceName+".conf")
	if err := os.Remove(configPath); err != nil {
		PrintErrorf("删除服务配置文件失败: %v", err)
		return err
	}

	if err := p.Reload(serviceName); err != nil {
		return fmt.Errorf("重载supervisor配置失败: %v", err)
	}

	return nil
}

// GetServiceList 获取所有服务列表
func (p *ServiceManager) GetServiceList() (string, error) {
	PrintInfo("获取服务列表...")

	if err := p.checkSupervisorInstalled(); err != nil {
		return "", err
	}

	output, err := RunShellWithOutput("supervisorctl", "status")
	if err != nil {
		return "", fmt.Errorf("获取服务列表失败: %v", err)
	}

	return string(output), nil
}

// --- 命令 ---

func (p *ServiceManager) GetServiceRootCommand() *cobra.Command {
	rootCmd := NewCommand(CommandOptions{
		Use:   "service",
		Short: "服务管理，管理系统中运行的后台服务，包括查看、启动、停止、重启等操作",
	})

	rootCmd.AddCommand(p.GetServiceListCommand())
	rootCmd.AddCommand(p.GetServiceStartCommand())
	rootCmd.AddCommand(p.GetServiceStopCommand())
	rootCmd.AddCommand(p.GetServiceRestartCommand())
	rootCmd.AddCommand(p.GetServiceStatusCommand())
	rootCmd.AddCommand(p.GetServiceLogsCommand())

	return rootCmd
}

// GetServiceListCommand 列出所有服务
func (p *ServiceManager) GetServiceListCommand() *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "list",
		Short: "列出所有服务",
		Run: func(cmd *cobra.Command, args []string) {
			// 检查 supervisor 是否正在运行
			if err := p.checkSupervisorRunning(); err != nil {
				return
			}

			PrintInfo("获取服务列表...")
			output, err := p.GetServiceList()
			if err != nil {
				PrintErrorf("%v", err)
				return
			}

			if output == "" {
				PrintInfo("当前没有运行中的服务")
				PrintCommandOutput("当前没有运行中的服务")
				return
			}

			PrintInfo("服务列表:")
			fmt.Println(output)
		},
	})
}

func (p *ServiceManager) GetServiceStartCommand() *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "start",
		Short: "启动服务",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceName := args[0]
			if err := p.Start(serviceName); err != nil {
				PrintErrorf("启动服务失败: %v", err)
				return
			}
		},
	})
}

func (p *ServiceManager) GetServiceStopCommand() *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "stop",
		Short: "停止服务",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceName := args[0]
			if err := p.Stop(serviceName); err != nil {
				PrintErrorf("停止服务失败: %v", err)
				return
			}
		},
	})
}

func (p *ServiceManager) GetServiceRestartCommand() *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "restart",
		Short: "重启服务",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceName := args[0]
			PrintInfof("正在重启服务: %s", serviceName)

			if err := p.Stop(serviceName); err != nil {
				PrintErrorf("停止服务失败: %v", err)
				return
			}

			if err := p.Start(serviceName); err != nil {
				PrintErrorf("启动服务失败: %v", err)
				return
			}

			PrintSuccessf("服务已成功重启: %s", serviceName)
		},
	})
}

func (p *ServiceManager) GetServiceStatusCommand() *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "status",
		Short: "查看服务状态",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceName := args[0]
			if err := p.checkSupervisorInstalled(); err != nil {
				return
			}

			command := exec.Command("supervisorctl", "status", serviceName)
			output, err := command.CombinedOutput()
			if err != nil {
				PrintErrorf("获取服务状态失败: %v", err)
				return
			}

			PrintInfof("服务状态:")
			fmt.Println(string(output))
		},
	})
}

func (p *ServiceManager) GetServiceLogsCommand() *cobra.Command {
	var tail int
	cmd := NewCommand(CommandOptions{
		Use:   "logs",
		Short: "查看服务日志",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceName := args[0]

			// 构建日志文件路径
			stdoutLog := filepath.Join(p.RootFolder, "logs", serviceName+".out.log")
			stderrLog := filepath.Join(p.RootFolder, "logs", serviceName+".err.log")

			// 读取标准输出日志
			PrintInfof("标准输出日志 (最后 %d 行):", tail)
			if err := p.tailLog(stdoutLog, tail); err != nil {
				PrintErrorf("读取标准输出日志失败: %v", err)
			}

			// 读取错误日志
			PrintInfof("\n错误日志 (最后 %d 行):", tail)
			if err := p.tailLog(stderrLog, tail); err != nil {
				PrintErrorf("读取错误日志失败: %v", err)
			}
		},
	})

	cmd.Flags().IntVarP(&tail, "tail", "n", 100, "显示最后几行日志")
	return cmd
}

func (p *ServiceManager) tailLog(logFile string, lines int) error {
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		PrintInfo("日志文件不存在")
		return nil
	}

	command := exec.Command("tail", "-n", fmt.Sprintf("%d", lines), logFile)
	output, err := command.CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
