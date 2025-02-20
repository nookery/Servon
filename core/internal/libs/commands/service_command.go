package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"servon/core/internal/libs/managers"

	"github.com/spf13/cobra"
)

// 将 ServiceManager 的能力通过命令行暴露出来

func GetServiceRootCommand(m *managers.ServiceManager) *cobra.Command {
	rootCmd := NewCommand(CommandOptions{
		Use:   "service",
		Short: "服务管理，管理系统中运行的后台服务，包括查看、启动、停止、重启等操作",
	})

	rootCmd.AddCommand(GetServiceListCommand(m))
	rootCmd.AddCommand(GetServiceStartCommand(m))
	rootCmd.AddCommand(GetServiceStopCommand(m))
	rootCmd.AddCommand(GetServiceRestartCommand(m))
	rootCmd.AddCommand(GetServiceStatusCommand(m))
	rootCmd.AddCommand(GetServiceLogsCommand(m))

	return rootCmd
}

// GetServiceListCommand 列出所有服务
func GetServiceListCommand(m *managers.ServiceManager) *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "list",
		Short: "列出所有服务",
		Run: func(cmd *cobra.Command, args []string) {
			// 检查 supervisor 是否正在运行
			if err := m.CheckSupervisorRunning(); err != nil {
				return
			}

			PrintInfo("获取服务列表...")
			output, err := m.GetServiceList()
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

func GetServiceStartCommand(m *managers.ServiceManager) *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "start",
		Short: "启动服务",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceName := args[0]
			if err := m.Start(serviceName); err != nil {
				PrintErrorf("启动服务失败: %v", err)
				return
			}
		},
	})
}

func GetServiceStopCommand(m *managers.ServiceManager) *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "stop",
		Short: "停止服务",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceName := args[0]
			if err := m.Stop(serviceName); err != nil {
				PrintErrorf("停止服务失败: %v", err)
				return
			}
		},
	})
}

func GetServiceRestartCommand(m *managers.ServiceManager) *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "restart",
		Short: "重启服务",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceName := args[0]
			PrintInfof("正在重启服务: %s", serviceName)

			if err := m.Stop(serviceName); err != nil {
				PrintErrorf("停止服务失败: %v", err)
				return
			}

			if err := m.Start(serviceName); err != nil {
				PrintErrorf("启动服务失败: %v", err)
				return
			}

			PrintSuccessf("服务已成功重启: %s", serviceName)
		},
	})
}

func GetServiceStatusCommand(m *managers.ServiceManager) *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "status",
		Short: "查看服务状态",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceName := args[0]
			if err := m.CheckSupervisorInstalled(); err != nil {
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

func GetServiceLogsCommand(m *managers.ServiceManager) *cobra.Command {
	var tail int
	cmd := NewCommand(CommandOptions{
		Use:   "logs",
		Short: "查看服务日志",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceName := args[0]

			// 构建日志文件路径
			stdoutLog := filepath.Join(m.RootFolder, "logs", serviceName+".out.log")
			stderrLog := filepath.Join(m.RootFolder, "logs", serviceName+".err.log")

			// 读取标准输出日志
			PrintInfof("标准输出日志 (最后 %d 行):", tail)
			if err := tailLog(stdoutLog, tail); err != nil {
				PrintErrorf("读取标准输出日志失败: %v", err)
			}

			// 读取错误日志
			PrintInfof("\n错误日志 (最后 %d 行):", tail)
			if err := tailLog(stderrLog, tail); err != nil {
				PrintErrorf("读取错误日志失败: %v", err)
			}
		},
	})

	cmd.Flags().IntVarP(&tail, "tail", "n", 100, "显示最后几行日志")
	return cmd
}

func tailLog(logFile string, lines int) error {
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
