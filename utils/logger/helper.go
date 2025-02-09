package logger

import (
	"fmt"
	"os/exec"
)

// 为方便使用，提供包级别的函数
func Debug(format string, args ...interface{}) {
	GetLogger().Debug(format, args...)
}

func DebugChan(ch chan<- string, format string, args ...interface{}) {
	GetLogger().DebugChan(ch, format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Info(format, args...)
}

func InfoWithSpace(format string, args ...interface{}) {
	fmt.Println()
	GetLogger().Info(format, args...)
	fmt.Println()
}

// InfoTitle 打印醒目的标题信息
// 用于在日志中突出显示重要的分段或章节标题
// 格式：=== 标题内容 ===
func InfoTitle(format string, args ...interface{}) {
	fmt.Println()
	fmt.Printf(colorBold+colorMagenta+"=== "+format+" ==="+colorReset+"\n", args...)
	fmt.Println()
}

func Warn(format string, args ...interface{}) {
	GetLogger().Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	GetLogger().Error(format, args...)
}

func InfoChan(ch chan<- string, format string, args ...interface{}) {
	GetLogger().InfoChan(ch, format, args...)
}

func ErrorChan(ch chan<- string, format string, args ...interface{}) {
	GetLogger().ErrorChan(ch, format, args...)
}

// StreamCommand 实时处理命令的输出流
func StreamCommand(cmd *exec.Cmd) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("无法创建标准输出管道: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("无法创建错误输出管道: %v", err)
	}

	Info("🚀 启动命令: %s", cmd.String())
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动命令失败: %v", err)
	}

	// 处理标准输出
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				fmt.Printf("%s", string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	// 处理错误输出
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				fmt.Printf("%s", string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	return cmd.Wait()
}
