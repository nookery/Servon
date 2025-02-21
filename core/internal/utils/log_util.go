package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LogLevel 定义日志级别
type LogLevel string

const (
	LogLevelInfo    LogLevel = "信息"
	LogLevelError   LogLevel = "错误"
	LogLevelWarning LogLevel = "警告"
	LogLevelDebug   LogLevel = "调试"
)

// LogUtil 提供日志记录功能
type LogUtil struct {
	logDir string
}

// NewLogUtil 创建新的日志工具实例
func NewLogUtil(logDir string) (*LogUtil, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %v", err)
	}
	return &LogUtil{logDir: logDir}, nil
}

// Close 关闭日志文件
func (l *LogUtil) Close() error {
	return nil
}

// Log 记录日志到主日志文件
func (l *LogUtil) Log(level LogLevel, format string, args ...interface{}) {
	// This method is now empty as the LogUtil struct no longer has a logger or file
}

// LogToFile 同时记录到主日志和指定的日志文件
func (l *LogUtil) LogToFile(level LogLevel, logFile *os.File, format string, args ...interface{}) {
	if logFile == nil {
		return
	}
	message := fmt.Sprintf(format, args...)
	timestamp := time.Now().In(time.Local).Format("2006-01-02 15:04:05.000")
	fmt.Fprintf(logFile, "[%s][%s] %s\n", timestamp, level, message)
}

// CreateLogFile 创建新的日志文件
func (l *LogUtil) CreateLogFile(name, header string) (string, *os.File, error) {
	timestamp := time.Now().In(time.Local)
	logID := timestamp.Format("2006-01-02-150405")
	logPath := filepath.Join(l.logDir, logID+".log")

	file, err := os.Create(logPath)
	if err != nil {
		return "", nil, fmt.Errorf("创建日志文件失败: %v", err)
	}

	// 写入日志头部信息，包含时间戳
	headerWithTime := fmt.Sprintf("[%s] === 日志开始 ===\n%s\n",
		timestamp.Format("2006-01-02 15:04:05.000"),
		header)

	if _, err = file.WriteString(headerWithTime); err != nil {
		file.Close()
		return "", nil, fmt.Errorf("写入日志头部信息失败: %v", err)
	}

	return logID, file, nil
}
