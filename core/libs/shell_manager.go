package libs

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

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
	go processOutput(stdoutPipe, "stdout")

	// 处理标准错误
	go processOutput(stderrPipe, "stderr")

	// 等待命令完成
	return execCmd.Wait()
}

// processOutput 处理输出流
func processOutput(pipe io.ReadCloser, source string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		// 打印到控制台并发送到日志通道
		fmt.Println(line)
		DefaultPrinter.sendToChannel(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取%s错误: %v\n", source, err)
	}
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
