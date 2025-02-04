package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var levelNames = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

// ANSI 颜色代码
const (
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[36m"
	colorGreen  = "\033[32m"
	colorReset  = "\033[0m"
)

var levelColors = map[LogLevel]string{
	DEBUG: colorBlue,
	INFO:  colorGreen,
	WARN:  colorYellow,
	ERROR: colorRed,
}

type Logger struct {
	mu       sync.Mutex
	file     *os.File
	filename string
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// GetLogger 获取默认日志记录器
func GetLogger() *Logger {
	once.Do(func() {
		var err error
		defaultLogger, err = NewLogger("servon")
		if err != nil {
			panic(fmt.Sprintf("初始化日志记录器失败: %v", err))
		}
	})
	return defaultLogger
}

// NewLogger 创建新的日志记录器
func NewLogger(name string) (*Logger, error) {
	// 确保日志目录存在
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 生成日志文件名
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", name, time.Now().Format("2006-01-02")))

	// 打开日志文件
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %v", err)
	}

	return &Logger{
		file:     file,
		filename: filename,
	}, nil
}

// rotateFile 检查并轮转日志文件
func (l *Logger) rotateFile() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 检查当前日期
	currentDate := time.Now().Format("2006-01-02")
	expectedFilename := filepath.Join("logs", fmt.Sprintf("servon_%s.log", currentDate))

	// 如果日期变化，创建新文件
	if l.filename != expectedFilename {
		// 关闭旧文件
		if err := l.file.Close(); err != nil {
			return fmt.Errorf("关闭旧日志文件失败: %v", err)
		}

		// 打开新文件
		file, err := os.OpenFile(expectedFilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("打开新日志文件失败: %v", err)
		}

		l.file = file
		l.filename = expectedFilename
	}

	return nil
}

// log 记录日志并可选择性地发送到channel
func (l *Logger) log(level LogLevel, ch chan<- string, format string, args ...interface{}) {
	// 检查并轮转日志文件
	if err := l.rotateFile(); err != nil {
		fmt.Printf("轮转日志文件失败: %v\n", err)
		return
	}

	// 获取调用信息
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "unknown"
		line = 0
	}
	// 查找并保留从 "internal" 开始的路径
	if idx := strings.Index(file, "internal"); idx != -1 {
		file = file[idx:]
	}

	// 格式化日志消息
	now := time.Now()
	// 文件日志使用完整时间格式
	fullTimeStr := now.Format("2006-01-02 15:04:05.000")
	// 终端输出使用简短时间格式
	shortTimeStr := now.Format("15:04:05")

	msg := fmt.Sprintf(format, args...)

	// 写入日志文件
	logLine := fmt.Sprintf("[%s] [%s] [%s:%d] %s\n",
		fullTimeStr, levelNames[level], file, line, msg)
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, err := l.file.WriteString(logLine); err != nil {
		fmt.Printf("写入日志失败: %v\n", err)
	}

	// 如果提供了channel，发送消息并跳过控制台输出
	if ch != nil {
		if level == ERROR {
			ch <- "Error: " + msg
		} else {
			ch <- msg
		}
		return
	}

	// 只有在没有提供channel时才输出到控制台
	coloredLogLine := fmt.Sprintf("%s[%s] [%s] [%s:%d] %s%s\n",
		levelColors[level], shortTimeStr, levelNames[level], file, line, msg, colorReset)
	fmt.Print(coloredLogLine)
}

// Debug 记录调试级别日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, nil, format, args...)
}

// Info 记录信息级别日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, nil, format, args...)
}

// Warn 记录警告级别日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, nil, format, args...)
}

// Error 记录错误级别日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, nil, format, args...)
}

// InfoChan 添加新的方法支持channel
func (l *Logger) InfoChan(ch chan<- string, format string, args ...interface{}) {
	l.log(INFO, ch, format, args...)
}

// ErrorChan 添加错误级别的channel支持
func (l *Logger) ErrorChan(ch chan<- string, format string, args ...interface{}) {
	l.log(ERROR, ch, format, args...)
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.file.Close()
}

// 为方便使用，提供包级别的函数
func Debug(format string, args ...interface{}) {
	GetLogger().Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Info(format, args...)
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
