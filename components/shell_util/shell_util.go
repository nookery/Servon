// Package shell_util 提供 Shell 命令执行和交互功能
//
// 这个组件封装了与 Shell 相关的工具函数，
// 包括命令执行、用户交互、颜色输出等功能。
package shell_util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

var DefaultShellUtil = ShellUtil{}

type ShellUtil struct{}

func NewShellUtil() *ShellUtil {
	return &ShellUtil{}
}

// ExecuteCommand 执行命令并返回输出
func (s *ShellUtil) ExecuteCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// ExecuteCommandWithOutput 执行命令并实时输出
func (s *ShellUtil) ExecuteCommandWithOutput(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = color.Output
	cmd.Stderr = color.Error
	return cmd.Run()
}

// ExecuteCommandInDir 在指定目录执行命令
func (s *ShellUtil) ExecuteCommandInDir(dir, command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// GetCurrentUser 获取当前用户信息
func (s *ShellUtil) GetCurrentUser() (*user.User, error) {
	return user.Current()
}

// IsCommandAvailable 检查命令是否可用
func (s *ShellUtil) IsCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// ReadUserInput 读取用户输入
func (s *ShellUtil) ReadUserInput(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// PrintSuccess 打印成功信息
func (s *ShellUtil) PrintSuccess(message string) {
	color.Green("✅ %s", message)
}

// PrintError 打印错误信息
func (s *ShellUtil) PrintError(message string) {
	color.Red("❌ %s", message)
}

// PrintWarning 打印警告信息
func (s *ShellUtil) PrintWarning(message string) {
	color.Yellow("⚠️  %s", message)
}

// PrintInfo 打印信息
func (s *ShellUtil) PrintInfo(message string) {
	color.Blue("ℹ️  %s", message)
}

// SplitCommand 分割命令字符串
func (s *ShellUtil) SplitCommand(command string) []string {
	return strings.Fields(command)
}

// JoinCommand 连接命令参数
func (s *ShellUtil) JoinCommand(args []string) string {
	return strings.Join(args, " ")
}

// EscapeShellArg 转义 Shell 参数
func (s *ShellUtil) EscapeShellArg(arg string) string {
	if strings.Contains(arg, " ") || strings.Contains(arg, "\t") {
		return fmt.Sprintf("'%s'", strings.ReplaceAll(arg, "'", "'\"'\"'"))
	}
	return arg
}

// GetShellType 获取当前 Shell 类型
func (s *ShellUtil) GetShellType() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "unknown"
	}
	return filepath.Base(shell)
}

// SetEnvironmentVariable 设置环境变量
func (s *ShellUtil) SetEnvironmentVariable(key, value string) error {
	return os.Setenv(key, value)
}

// GetEnvironmentVariable 获取环境变量
func (s *ShellUtil) GetEnvironmentVariable(key string) string {
	return os.Getenv(key)
}

// ExecuteCommandAsync 异步执行命令
func (s *ShellUtil) ExecuteCommandAsync(command string) (*exec.Cmd, error) {
	cmd := exec.Command("sh", "-c", command)
	err := cmd.Start()
	return cmd, err
}

// WaitForCommand 等待命令完成
func (s *ShellUtil) WaitForCommand(cmd *exec.Cmd) error {
	return cmd.Wait()
}

// KillCommand 终止命令
func (s *ShellUtil) KillCommand(cmd *exec.Cmd) error {
	if cmd.Process != nil {
		return cmd.Process.Kill()
	}
	return fmt.Errorf("no process to kill")
}

// RunShell 执行命令
func (c *ShellUtil) RunShell(command string, args ...string) (error, string) {
	return c.execute("", false, command, args...)
}

// RunShellWithOutput 执行命令并返回输出
func (c *ShellUtil) RunShellWithOutput(command string, args ...string) (error, string) {
	return c.execute("", false, command, args...)
}

// RunShellWithSudo 使用sudo执行命令
func (c *ShellUtil) RunShellWithSudo(command string, args ...string) (error, string) {
	return c.execute("", true, command, args...)
}

// RunShellInFolder 在指定目录中执行命令
func (c *ShellUtil) RunShellInFolder(dir string, command string, args ...string) (error, string) {
	return c.execute(dir, false, command, args...)
}

// RunShellWithSudoInFolder 在指定目录中使用sudo执行命令
func (c *ShellUtil) RunShellWithSudoInFolder(dir string, command string, args ...string) (error, string) {
	return c.execute(dir, true, command, args...)
}

// RunShellWithSudoOutput 使用sudo执行命令并返回输出
func (c *ShellUtil) RunShellWithSudoOutput(command string, args ...string) (error, string) {
	return c.execute("", true, command, args...)
}

// isRoot 检查是否为root用户
func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		return false
	}
	return currentUser.Uid == "0"
}

// execute 是内部函数，负责实际的命令执行逻辑
func (c *ShellUtil) execute(dir string, withSudo bool, command string, args ...string) (error, string) {
	if command == "" {
		return fmt.Errorf("command is required"), ""
	}

	// 构建完整的命令参数
	var cmdArgs []string

	// 处理 sh -c "复杂命令" 的特殊情况
	if command == "sh" && len(args) >= 2 && args[0] == "-c" {
		// 如果需要 sudo 并且不是 root 用户
		if withSudo && !isRoot() {
			cmdArgs = append(cmdArgs, "sudo", "sh", "-c", args[1])
		} else {
			cmdArgs = append(cmdArgs, "sh", "-c", args[1])
		}
	} else {
		// 常规命令处理
		if withSudo && !isRoot() {
			cmdArgs = append(cmdArgs, "sudo")
		}
		cmdArgs = append(cmdArgs, command)
		cmdArgs = append(cmdArgs, args...)
	}

	// 创建命令对象
	var cmd *exec.Cmd
	if command == "sh" && len(args) >= 2 && args[0] == "-c" {
		// 对于 sh -c "复杂命令"，我们需要特殊处理
		if withSudo && !isRoot() {
			cmd = exec.Command("sudo", "sh", "-c", args[1])
		} else {
			cmd = exec.Command("sh", "-c", args[1])
		}
	} else {
		cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
	}

	// 设置工作目录（如果指定）
	if dir != "" {
		cmd.Dir = dir
	}

	// 创建输出缓冲区
	var output strings.Builder

	// 设置标准输出和错误输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err), ""
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err), ""
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err), ""
	}

	// 读取输出
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			output.WriteString(line + "\n")
			color.Green(line)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			output.WriteString(line + "\n")
			color.Red(line)
		}
	}()

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command execution failed: %v", err), output.String()
	}

	return nil, output.String()
}