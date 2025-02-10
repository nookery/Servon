package api

import (
	"servon/core/libs"
)

type LogApi struct{}

func NewLogApi() LogApi {
	return LogApi{}
}

// Error 打印错误信息
func (c *LogApi) Error(args ...interface{}) {
	libs.Error(args...)
}

// Errorf 打印错误信息
func (c *LogApi) Errorf(format string, args ...interface{}) {
	libs.Errorf(format, args...)
}

// ErrorChan 打印错误信息到通道
func (c *LogApi) ErrorChan(ch chan<- string, format string, args ...interface{}) {
	libs.ErrorChan(ch, format, args...)
}

// Info 打印信息
func (c *LogApi) Info(format string, args ...interface{}) {
	libs.Info(format, args...)
}

// InfoWithSpace 打印信息并换行
func (c *LogApi) InfoWithSpace(format string, args ...interface{}) {
	libs.InfoWithSpace(format, args...)
}

// Infoln 打印信息并换行
func (c *LogApi) Infoln(format string, args ...interface{}) {
	libs.Infoln(format, args...)
}

// InfoChan 打印信息到通道
func (c *LogApi) InfoChan(ch chan<- string, format string, args ...interface{}) {
	libs.InfoChan(ch, format, args...)
}

func (c *LogApi) PrintAndReturnError(errMsg string) error {
	return libs.PrintAndReturnError(errMsg)
}

// PrintList 打印列表
func (c *LogApi) PrintList(list []string, title string) {
	libs.PrintList(list, title)
}

// Warn 打印警告信息
func (c *LogApi) Warn(format string, args ...interface{}) {
	libs.Warn(format, args...)
}
