package logger

import "github.com/fatih/color"

// LogType 定义日志类型及其属性
type LogType struct {
	Name   string
	Color  *color.Color
	Symbol string
}

// 定义所有日志类型
var (
	LogTypeInfo = LogType{
		Name:   "info",
		Color:  color.New(color.FgCyan),
		Symbol: "🍋",
	}
	LogTypeError = LogType{
		Name:   "error",
		Color:  color.New(color.FgRed),
		Symbol: "❌",
	}
	LogTypeWarn = LogType{
		Name:   "warn",
		Color:  color.New(color.FgYellow),
		Symbol: "🚨",
	}
	LogTypeSuccess = LogType{
		Name:   "success",
		Color:  color.New(color.FgGreen),
		Symbol: "✅",
	}
	LogTypeDebug = LogType{
		Name:   "debug",
		Color:  color.New(color.FgBlue),
		Symbol: "🔍",
	}
)
