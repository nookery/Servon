package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

var DefaultLogUtil = NewConsoleUtil()

type LogUtil struct {
	logDir string
	logger zerolog.Logger
}

// NewConsoleUtil 初始化记录到控制台的日志工具
func NewConsoleUtil() *LogUtil {
	return newLogUtil("")
}

// NewLogUtil 初始化记录到文件和控制台的日志工具
func NewLogUtil(logDir string) *LogUtil {
	return newLogUtil(logDir)
}

func newLogUtil(logDir string) *LogUtil {
	var multi zerolog.LevelWriter

	if logDir != "" {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			panic(err)
		}

		// 创建日志文件
		logFile := filepath.Join(logDir, "app.log")
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		// 同时输出到控制台和文件
		multi = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, file)
	} else {
		// 仅输出到控制台
		multi = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	// 配置 zerolog
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	logger := zerolog.New(multi).With().Timestamp().
		CallerWithSkipFrameCount(3). // 跳过日志工具的调用栈
		Logger()

	return &LogUtil{
		logDir: logDir,
		logger: logger,
	}
}

// Alert 记录警告日志
func (lu *LogUtil) Alert(message string) {
	lu.logger.Warn().Msg(message)
}

// Alertf 记录警告日志
func (lu *LogUtil) Alertf(format string, args ...interface{}) {
	lu.logger.Warn().Msgf(format, args...)
}

func (lu *LogUtil) Info(message string) {
	lu.logger.Info().Msg(message)
}

func (lu *LogUtil) Infof(format string, args ...interface{}) {
	lu.logger.Info().Msgf(format, args...)
}

func (lu *LogUtil) Trace(message string) {
	lu.logger.Trace().Msg(message)
}

func (lu *LogUtil) Error(err error) {
	lu.logger.Error().Msg(err.Error())
}

func (lu *LogUtil) ErrorMessage(message string) {
	lu.logger.Error().Msg(message)
}

func (lu *LogUtil) Errorf(format string, args ...interface{}) {
	lu.logger.Error().Msgf(format, args...)
}

// LogAndReturnError 记录错误日志并返回错误
func (lu *LogUtil) LogAndReturnError(message string) error {
	lu.ErrorMessage(message)
	return fmt.Errorf("%s", message)
}

// LogAndReturnErrorf 记录错误日志并返回错误
func (lu *LogUtil) LogAndReturnErrorf(format string, args ...interface{}) error {
	lu.Errorf(format, args...)
	return fmt.Errorf(format, args...)
}

func (lu *LogUtil) Debug(message string) {
	lu.logger.Debug().Msg(message)
}

func (lu *LogUtil) Debugf(format string, args ...interface{}) {
	lu.logger.Debug().Msgf(format, args...)
}

func (lu *LogUtil) Warn(message string) {
	lu.logger.Warn().Msg(message)
}

func (lu *LogUtil) Warnf(format string, args ...interface{}) {
	lu.logger.Warn().Msgf(format, args...)
}

func (lu *LogUtil) Fatal(message string) {
	lu.logger.Fatal().Msg(message)
}

// Success 记录成功日志
func (lu *LogUtil) Success(message string) {
	lu.logger.Info().Msg(message)
}

// Successf 记录成功日志
func (lu *LogUtil) Successf(format string, args ...interface{}) {
	lu.logger.Info().Msgf(format, args...)
}

// EmptyLine 记录空行
func (lu *LogUtil) EmptyLine() {
	lu.logger.Info().Msg("")
}

// Title 记录标题日志
func (lu *LogUtil) Title(message string) {
	lu.logger.Info().Msg(message)
}

// Titlef 记录标题日志
func (lu *LogUtil) Titlef(format string, args ...interface{}) {
	lu.logger.Info().Msgf(format, args...)
}

// List 记录列表日志
func (lu *LogUtil) List(message string) {
	lu.logger.Info().Msg(message)
}

// ListWithTitle 记录列表日志
func (lu *LogUtil) ListWithTitle(title string, list []string) {
	lu.logger.Info().Msgf("%s: %s", title, strings.Join(list, ", "))
}

// Listf 记录列表日志
func (lu *LogUtil) Listf(format string, args ...interface{}) {
	lu.logger.Info().Msgf(format, args...)
}

// PrintKeyValue 记录键值对日志
func (lu *LogUtil) PrintKeyValue(key string, value string) {
	lu.logger.Info().Msgf("%s: %s", key, value)
}

// PrintKeyValues 记录键值对日志
func (lu *LogUtil) PrintKeyValues(keyValues map[string]string) {
	for key, value := range keyValues {
		lu.logger.Info().Msgf("%s: %s", key, value)
	}
}
