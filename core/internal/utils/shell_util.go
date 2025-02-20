package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

var DefaultShellUtil = ShellUtil{}

type ShellUtil struct{}

func NewShellUtil() *ShellUtil {
	return &ShellUtil{}
}

func (c *ShellUtil) Execute(command string, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("command is required")
	}

	// 使用青色（Cyan）输出命令和参数，用空格连接参数
	color.Cyan("📺 %s %s", command, JoinArgs(args))

	execCmd := exec.Command(command, args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	return execCmd.Run()
}

// StreamCommand 执行命令并打印输出
func (c *ShellUtil) StreamCommand(cmd *exec.Cmd) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("获取标准输出失败: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("获取标准错误输出失败: %v", err)
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	return cmd.Run()
}

func (c *ShellUtil) ExecuteWithOutput(command string, args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("command is required")
	}

	execCmd := exec.Command(command, args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	output, err := execCmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
