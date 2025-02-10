package libs

import (
	"bufio"
	"fmt"
	"os/exec"
)

func Debug(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func Error(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func ErrorChan(ch chan<- string, format string, args ...interface{}) {
	ch <- fmt.Sprintf(format, args...)
}

// Info 打印信息
func Info(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Infoln 打印信息并换行
func Infoln(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// InfoWithSpace 打印信息并换行
func InfoWithSpace(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// InfoChan 打印信息到通道
func InfoChan(ch chan<- string, format string, args ...interface{}) {
	Info(fmt.Sprintf(format, args...))
	ch <- fmt.Sprintf(format, args...)
}

func PrintAndReturnError(errMsg string) error {
	return fmt.Errorf(errMsg)
}

func PrintCommandErrorAndExit(err error) error {
	return fmt.Errorf(err.Error())
}

// StreamCommand 执行命令并打印输出
func StreamCommand(cmd *exec.Cmd) error {
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
