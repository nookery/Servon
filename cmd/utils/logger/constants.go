package logger

// LogLevel 定义日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var levelNames = map[LogLevel]string{
	DEBUG: "D",
	INFO:  "I",
	WARN:  "W",
	ERROR: "E",
}

// ANSI 颜色代码
const (
	colorRed     = "\033[31m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[36m"
	colorGreen   = "\033[32m"
	colorReset   = "\033[0m"
	colorMagenta = "\033[35m"
	colorBold    = "\033[1m"
)

var levelColors = map[LogLevel]string{
	DEBUG: colorBlue,
	INFO:  colorGreen,
	WARN:  colorYellow,
	ERROR: colorRed,
}
