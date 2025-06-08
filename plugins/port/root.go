package port

import (
	"fmt"
	"os/exec"
	"runtime"
	"servon/core"
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// PortCmd represents the port command
var PortCmd = &cobra.Command{
	Use:   "port",
	Short: color.Blue.Render("端口管理工具"),
	Long:  color.Success.Render("\r\n端口管理工具，用于查看和关闭占用指定端口的程序"),
}

func init() {
	// 添加子命令
	PortCmd.AddCommand(killCmd)
	PortCmd.AddCommand(listCmd)
}

// killCmd 关闭占用指定端口的程序
var killCmd = &cobra.Command{
	Use:   "kill [port]",
	Short: "关闭占用指定端口的程序",
	Long:  color.Success.Render("\r\n关闭占用指定端口的程序，支持强制关闭选项"),
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		port := args[0]
		force, _ := cmd.Flags().GetBool("force")
		verbose, _ := cmd.Flags().GetBool("verbose")

		// 验证端口号
		portNum, err := strconv.Atoi(port)
		if err != nil || portNum < 1 || portNum > 65535 {
			return fmt.Errorf("无效的端口号: %s，端口号必须在 1-65535 之间", port)
		}

		color.Info.Printf("🔍 正在查找占用端口 %s 的程序...\n", port)

		// 显示使用的检测命令
		var detectionCmd string
		switch runtime.GOOS {
		case "darwin":
			detectionCmd = fmt.Sprintf("lsof -i :%s", port)
		case "linux":
			detectionCmd = fmt.Sprintf("lsof -i :%s", port)
		case "windows":
			detectionCmd = "netstat -ano"
		}
		color.Info.Printf("💡 检测命令: %s\n", detectionCmd)

		// 查找占用端口的进程
		pids, err := findProcessByPort(port)
		if err != nil {
			return fmt.Errorf("查找进程失败: %v", err)
		}

		if len(pids) == 0 {
			color.Yellow.Printf("⚠️  端口 %s 未被占用\n", port)
			return nil
		}

		color.Info.Printf("📋 找到 %d 个占用端口 %s 的进程:\n", len(pids), port)

		// 显示进程信息
		for _, pid := range pids {
			processInfo, err := getProcessInfo(pid)
			if err != nil {
				color.Error.Printf("❌ 获取进程 %s 信息失败: %v\n", pid, err)
				continue
			}
			color.Info.Printf("  PID: %s, 进程名: %s\n", pid, processInfo)
		}

		// 关闭进程
		successCount := 0
		for _, pid := range pids {
			// 检查是否为 unknown PID
			if pid == "unknown" {
				color.Yellow.Printf("⚠️  无法关闭进程：netstat 无法获取 PID，请手动查找并关闭占用端口 %s 的进程\n", port)
				continue
			}

			if verbose {
				color.Info.Printf("🔄 正在关闭进程 %s...\n", pid)
			}

			err := killProcess(pid, force)
			if err != nil {
				color.Error.Printf("❌ 关闭进程 %s 失败: %v\n", pid, err)
				continue
			}

			color.Success.Printf("✅ 成功关闭进程 %s\n", pid)
			successCount++
		}

		if successCount > 0 {
			color.Success.Printf("🎉 端口 %s 已释放\n", port)
		} else {
			color.Yellow.Printf("⚠️  端口 %s 的进程未能关闭，请手动处理\n", port)
		}
		return nil
	},
}

// listCmd 列出占用指定端口的程序
var listCmd = &cobra.Command{
	Use:   "list [port]",
	Short: "列出占用指定端口的程序",
	Long:  color.Success.Render("\r\n列出占用指定端口的程序信息，不进行关闭操作"),
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		port := args[0]
		verbose, _ := cmd.Flags().GetBool("verbose")

		// 验证端口号
		portNum, err := strconv.Atoi(port)
		if err != nil || portNum < 1 || portNum > 65535 {
			return fmt.Errorf("无效的端口号: %s，端口号必须在 1-65535 之间", port)
		}

		color.Info.Printf("🔍 正在查找占用端口 %s 的程序...\n", port)

		// 显示使用的检测命令
		var detectionCmd string
		switch runtime.GOOS {
		case "darwin":
			detectionCmd = fmt.Sprintf("lsof -i :%s", port)
		case "linux":
			detectionCmd = fmt.Sprintf("lsof -i :%s", port)
		case "windows":
			detectionCmd = "netstat -ano"
		}
		color.Info.Printf("💡 检测命令: %s\n", detectionCmd)

		// 查找占用端口的进程
		pids, err := findProcessByPort(port)
		if err != nil {
			return fmt.Errorf("查找进程失败: %v", err)
		}

		if len(pids) == 0 {
			color.Yellow.Printf("⚠️  端口 %s 未被占用\n", port)
			return nil
		}

		color.Info.Printf("📋 找到 %d 个占用端口 %s 的进程:\n", len(pids), port)

		// 显示进程详细信息
		for _, pid := range pids {
			processInfo, err := getProcessInfo(pid)
			if err != nil {
				color.Error.Printf("❌ 获取进程 %s 信息失败: %v\n", pid, err)
				continue
			}

			color.Info.Printf("  PID: %s\n", pid)
			color.Info.Printf("  进程名: %s\n", processInfo)

			if verbose {
				// 获取更详细的进程信息
				detailedInfo, err := getDetailedProcessInfo(pid)
				if err == nil {
					color.Info.Printf("  详细信息: %s\n", detailedInfo)
				}
			}
			color.Info.Println("  ---")
		}

		return nil
	},
}

func init() {
	// 为 kill 命令添加标志
	killCmd.Flags().BoolP("force", "f", false, "强制关闭进程 (使用 SIGKILL)")
	killCmd.Flags().BoolP("verbose", "v", false, "显示详细信息")

	// 为 list 命令添加标志
	listCmd.Flags().BoolP("verbose", "v", false, "显示详细信息")
}

// findProcessByPort 根据端口号查找占用该端口的进程ID
func findProcessByPort(port string) ([]string, error) {
	var cmd *exec.Cmd
	var pids []string

	switch runtime.GOOS {
	case "darwin":
		// 在 macOS 上使用 lsof 命令查找占用端口的进程
		cmd = exec.Command("lsof", "-i", fmt.Sprintf(":%s", port))
	case "linux":
		// 使用 lsof 命令查找占用端口的进程
		cmd = exec.Command("lsof", "-i", fmt.Sprintf(":%s", port))
	case "windows":
		// 使用 netstat 命令查找占用端口的进程
		cmd = exec.Command("netstat", "-ano")
	default:
		return nil, fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		// 如果 lsof 命令失败，可能是因为没有找到进程
		if runtime.GOOS != "windows" {
			return pids, nil
		}
		return nil, fmt.Errorf("执行命令失败: %v", err)
	}

	switch runtime.GOOS {
	case "windows":
		// 解析 Windows netstat 输出
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, ":"+port+" ") {
				fields := strings.Fields(line)
				if len(fields) >= 5 {
					pid := fields[len(fields)-1]
					if pid != "0" {
						pids = append(pids, pid)
					}
				}
			}
		}
	case "darwin":
		// 解析 macOS lsof 输出
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for i, line := range lines {
			if i == 0 {
				// 跳过标题行
				continue
			}
			line = strings.TrimSpace(line)
			if line != "" {
				// 提取PID（第二列）
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					pid := fields[1]
					// 避免重复添加相同的PID
					found := false
					for _, existingPid := range pids {
						if existingPid == pid {
							found = true
							break
						}
					}
					if !found {
						pids = append(pids, pid)
					}
				}
			}
		}
	case "linux":
		// 解析 Linux lsof 输出
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for i, line := range lines {
			if i == 0 {
				// 跳过标题行
				continue
			}
			line = strings.TrimSpace(line)
			if line != "" {
				// 提取PID（第二列）
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					pid := fields[1]
					// 避免重复添加相同的PID
					found := false
					for _, existingPid := range pids {
						if existingPid == pid {
							found = true
							break
						}
					}
					if !found {
						pids = append(pids, pid)
					}
				}
			}
		}
	}

	return pids, nil
}

// getProcessInfo 获取进程基本信息
func getProcessInfo(pid string) (string, error) {
	// 处理 unknown PID 的情况
	if pid == "unknown" {
		return "未知进程 (netstat 无法获取 PID)", nil
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin", "linux":
		cmd = exec.Command("ps", "-p", pid, "-o", "comm=")
	case "windows":
		cmd = exec.Command("tasklist", "/fi", fmt.Sprintf("PID eq %s", pid), "/fo", "csv", "/nh")
	default:
		return "", fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取进程信息失败: %v", err)
	}

	processName := strings.TrimSpace(string(output))
	if runtime.GOOS == "windows" {
		// 解析 Windows tasklist CSV 输出
		fields := strings.Split(processName, ",")
		if len(fields) > 0 {
			processName = strings.Trim(fields[0], "\"")
		}
	}

	return processName, nil
}

// getDetailedProcessInfo 获取进程详细信息
func getDetailedProcessInfo(pid string) (string, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin", "linux":
		cmd = exec.Command("ps", "-p", pid, "-o", "pid,ppid,user,command")
	case "windows":
		cmd = exec.Command("wmic", "process", "where", fmt.Sprintf("ProcessId=%s", pid), "get", "Name,ParentProcessId,CommandLine", "/format:csv")
	default:
		return "", fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取详细进程信息失败: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// killProcess 关闭指定PID的进程
func killProcess(pid string, force bool) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin", "linux":
		if force {
			cmd = exec.Command("kill", "-9", pid)
		} else {
			cmd = exec.Command("kill", pid)
		}
	case "windows":
		if force {
			cmd = exec.Command("taskkill", "/F", "/PID", pid)
		} else {
			cmd = exec.Command("taskkill", "/PID", pid)
		}
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("关闭进程失败: %v", err)
	}

	return nil
}

// Setup 注册port插件到应用程序
func Setup(app *core.App) {
	app.GetRootCommand().AddCommand(PortCmd)
}