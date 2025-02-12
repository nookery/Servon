package libs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

type ShellManager struct {
}

func NewShellManager() *ShellManager {
	return &ShellManager{}
}

func (s *ShellManager) RunShell(command string, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("command is required")
	}

	// 使用青色（Cyan）输出命令和参数，用空格连接参数
	PrintInfo("📺 %s %s", command, joinArgs(args))

	execCmd := exec.Command(command, args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	return execCmd.Run()
}

// RunShellWithOutput 运行命令并返回输出
func (s *ShellManager) RunShellWithOutput(command string, args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("command is required")
	}

	color.Cyan("📺 %s %s", command, joinArgs(args))

	execCmd := exec.Command(command, args...)

	output, err := execCmd.CombinedOutput()

	Print(string(output))

	return string(output), err
}

// RunShellAndSendLog 运行命令并发送日志
func (s *ShellManager) RunShellAndSendLog(logChan chan<- string, command string, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("command is required")
	}

	PrintInfo("📺 %s %s", command, joinArgs(args))

	execCmd := exec.Command(command, args...)

	// 创建管道用于捕获输出
	stdoutPipe, err := execCmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderrPipe, err := execCmd.StderrPipe()
	if err != nil {
		return err
	}

	// 启动命令
	if err := execCmd.Start(); err != nil {
		return err
	}

	// 处理标准输出
	go func() {
		reader := bufio.NewReader(stdoutPipe)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Printf("读取stdout错误: %v\n", err)
				}
				break
			}
			line = strings.TrimRight(line, "\n")
			fmt.Println(line)
			logChan <- line
		}
	}()

	// 处理标准错误
	go func() {
		reader := bufio.NewReader(stderrPipe)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Printf("读取stderr错误: %v\n", err)
				}
				break
			}
			line = strings.TrimRight(line, "\n")
			fmt.Println(line)
			logChan <- line
		}
	}()

	// 等待命令完成
	return execCmd.Wait()
}
