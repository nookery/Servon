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
	logDir        string
	topic         string
	logger        zerolog.Logger
	consoleLogger zerolog.Logger
}

// NewConsoleUtil 初始化记录到控制台的日志工具
func NewConsoleUtil() *LogUtil {
	return newLogUtil("", "")
}

// NewLogUtil 初始化记录到文件和控制台的日志工具
func NewLogUtil(logDir string) *LogUtil {
	return newLogUtil(logDir, "")
}

// NewTopicLogUtil 初始化带主题的日志工具
func NewTopicLogUtil(logDir string, topic string) *LogUtil {
	return newLogUtil(logDir, topic)
}

func newLogUtil(logDir string, topic string) *LogUtil {
	var multi zerolog.LevelWriter
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}

	if logDir != "" {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			panic(err)
		}

		// 根据主题创建日志文件
		var logFile string
		if topic != "" {
			logFile = filepath.Join(logDir, fmt.Sprintf("%s.log", topic))
		} else {
			logFile = filepath.Join(logDir, "app.log")
		}

		// 使用自定义的文件写入器，支持自动重建文件
		file, err := newAutoCreateFile(logFile)
		if err != nil {
			panic(err)
		}
		multi = zerolog.MultiLevelWriter(consoleWriter, file)
	} else {
		multi = zerolog.MultiLevelWriter(consoleWriter)
	}

	// 配置 zerolog
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"

	// 创建基础logger配置
	baseLoggerContext := zerolog.New(multi).With().Timestamp().
		CallerWithSkipFrameCount(3)

	// 如果有主题，添加到日志上下文
	if topic != "" {
		baseLoggerContext = baseLoggerContext.Str("topic", topic)
	}

	// 创建标准logger
	logger := baseLoggerContext.Logger()

	// 创建仅控制台输出的logger
	consoleLoggerContext := zerolog.New(consoleWriter).With().Timestamp().
		CallerWithSkipFrameCount(3)
	if topic != "" {
		consoleLoggerContext = consoleLoggerContext.Str("topic", topic)
	}
	consoleLogger := consoleLoggerContext.Logger()

	return &LogUtil{
		logDir:        logDir,
		topic:         topic,
		logger:        logger,
		consoleLogger: consoleLogger,
	}
}

// autoCreateFile 是一个支持自动重建的文件写入器
type autoCreateFile struct {
	filename string
	file     *os.File
}

func newAutoCreateFile(filename string) (*autoCreateFile, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &autoCreateFile{
		filename: filename,
		file:     file,
	}, nil
}

func (f *autoCreateFile) Write(p []byte) (n int, err error) {
	if f.file == nil {
		// 如果文件不存在，尝试重新创建
		file, err := os.OpenFile(f.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return 0, err
		}
		f.file = file
	}

	n, err = f.file.Write(p)
	if err != nil {
		// 如果写入出错（可能是文件被删除），关闭当前文件句柄
		f.file.Close()
		f.file = nil
		return n, err
	}
	return n, nil
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
	lu.logger.Error().Msg("❌ " + err.Error())
}

func (lu *LogUtil) ErrorMessage(message string) {
	lu.logger.Error().Msg("❌ " + message)
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

// WarnfConsole 记录警告日志到控制台
func (lu *LogUtil) WarnfConsole(format string, args ...interface{}) {
	lu.consoleLogger.Warn().Msgf(format, args...)
}

func (lu *LogUtil) Fatal(message string) {
	lu.logger.Fatal().Msg(message)
}

// Success 记录成功日志
func (lu *LogUtil) Success(message string) {
	lu.logger.Info().Msg("✅ " + message)
}

// Successf 记录成功日志
func (lu *LogUtil) Successf(format string, args ...interface{}) {
	lu.logger.Info().Msgf("✅ "+format, args...)
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

// InfoConsole 记录信息日志到控制台
func (lu *LogUtil) InfoConsole(message string) {
	lu.consoleLogger.Info().Msg(message)
}

// InfofConsole 记录信息日志到控制台
func (lu *LogUtil) InfofConsole(format string, args ...interface{}) {
	lu.consoleLogger.Info().Msgf(format, args...)
}

// ErrorConsole 记录错误日志到控制台
func (lu *LogUtil) ErrorConsole(err error) {
	lu.consoleLogger.Error().Msg(err.Error())
}

// ErrorfConsole 记录错误日志到控制台
func (lu *LogUtil) ErrorfConsole(format string, args ...interface{}) {
	lu.consoleLogger.Error().Msgf(format, args...)
}
