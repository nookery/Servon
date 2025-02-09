package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Logger 日志记录器结构体
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

	currentDate := time.Now().Format("2006-01-02")
	expectedFilename := filepath.Join("logs", fmt.Sprintf("servon_%s.log", currentDate))

	if l.filename != expectedFilename {
		if err := l.file.Close(); err != nil {
			return fmt.Errorf("关闭旧日志文件失败: %v", err)
		}

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
	if err := l.rotateFile(); err != nil {
		fmt.Printf("轮转日志文件失败: %v\n", err)
		return
	}

	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "unknown"
		line = 0
	}

	now := time.Now()
	fullTimeStr := now.Format("2006-01-02 15:04:05.000")
	shortTimeStr := now.Format("15:04:05")

	msg := fmt.Sprintf(format, args...)

	logLine := fmt.Sprintf("[%s] [%s] [%s:%d] %s\n",
		fullTimeStr, levelNames[level], file, line, msg)
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, err := l.file.WriteString(logLine); err != nil {
		fmt.Printf("写入日志失败: %v\n", err)
	}

	if ch != nil {
		if level == ERROR {
			ch <- "Error: " + msg
		} else {
			ch <- msg
		}
		return
	}

	coloredLogLine := fmt.Sprintf("%s[%s] [%s] [%s:%d] %s%s\n",
		levelColors[level], shortTimeStr, levelNames[level], file, line, msg, colorReset)
	fmt.Print(coloredLogLine)
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.file.Close()
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

// DebugChan 添加调试级别的channel支持
func (l *Logger) DebugChan(ch chan<- string, format string, args ...interface{}) {
	l.log(DEBUG, ch, format, args...)
}

// ErrorChan 添加错误级别的channel支持
func (l *Logger) ErrorChan(ch chan<- string, format string, args ...interface{}) {
	l.log(ERROR, ch, format, args...)
}
