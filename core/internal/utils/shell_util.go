package utils

import (
	"bufio"
	"fmt"
	"os/exec"
	"os/user"
	"strings"

	"github.com/fatih/color"
)

var DefaultShellUtil = ShellUtil{}

type ShellUtil struct{}

func NewShellUtil() *ShellUtil {
	return &ShellUtil{}
}

// RunShell 执行命令
func (c *ShellUtil) RunShell(command string, args ...string) (error, string) {
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

// RunShellWithOutput 执行命令并返回输出
func (c *ShellUtil) RunShellWithOutput(command string, args ...string) (error, string) {
	return c.execute("", false, command, args...)
}

// RunShellWithSudoOutput 使用sudo执行命令并返回输出
func (c *ShellUtil) RunShellWithSudoOutput(command string, args ...string) (error, string) {
	return c.execute("", true, command, args...)
}

// execute 是内部函数，负责实际的命令执行逻辑
func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		return false
	}
	return currentUser.Uid == "0"
}

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
