package utils

import (
	"fmt"
	"os"
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

// log 记录日志
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
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
	logLine := fmt.Sprintf("[%s] [%s] [%s:%d] %s\n",
		fullTimeStr, levelNames[level], file, line, msg)

	// 写入日志文件（不带颜色，使用完整时间）
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, err := l.file.WriteString(logLine); err != nil {
		fmt.Printf("写入日志失败: %v\n", err)
	}

	// 输出到控制台（带颜色，使用简短时间）
	coloredLogLine := fmt.Sprintf("%s[%s] [%s] [%s:%d] %s%s\n",
		levelColors[level], shortTimeStr, levelNames[level], file, line, msg, colorReset)
	fmt.Print(coloredLogLine)
}

// Debug 记录调试级别日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info 记录信息级别日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn 记录警告级别日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error 记录错误级别日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
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
