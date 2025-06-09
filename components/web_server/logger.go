package web_server

import (
	"log"
	"os"
)

// Logger 定义日志接口
type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

// DefaultLogger 默认日志实现
type DefaultLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	warnLogger  *log.Logger
	debugLogger *log.Logger
}

// NewDefaultLogger 创建默认日志器
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
		warnLogger:  log.New(os.Stdout, "[WARN] ", log.LstdFlags),
		debugLogger: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags),
	}
}

// Infof 输出信息日志
func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	l.infoLogger.Printf(format, args...)
}

// Errorf 输出错误日志
func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	l.errorLogger.Printf(format, args...)
}

// Warnf 输出警告日志
func (l *DefaultLogger) Warnf(format string, args ...interface{}) {
	l.warnLogger.Printf(format, args...)
}

// Debugf 输出调试日志
func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	l.debugLogger.Printf(format, args...)
}

// NoOpLogger 空日志实现（不输出任何日志）
type NoOpLogger struct{}

// NewNoOpLogger 创建空日志器
func NewNoOpLogger() *NoOpLogger {
	return &NoOpLogger{}
}

// Infof 空实现
func (l *NoOpLogger) Infof(format string, args ...interface{}) {}

// Errorf 空实现
func (l *NoOpLogger) Errorf(format string, args ...interface{}) {}

// Warnf 空实现
func (l *NoOpLogger) Warnf(format string, args ...interface{}) {}

// Debugf 空实现
func (l *NoOpLogger) Debugf(format string, args ...interface{}) {}
