package api

import (
	"servon/core/libs"
)

type LogApi struct{}

func NewLogApi() LogApi {
	return LogApi{}
}

// Error 打印错误信息
func (c *LogApi) Error(format string, args ...interface{}) {
	libs.Error(format, args...)
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
