package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

type LogUtil struct {
	logDir string
	logger zerolog.Logger
}

func NewLogUtil(logDir string) *LogUtil {
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
	multi := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, file)

	// 配置 zerolog
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	logger := zerolog.New(multi).With().Timestamp().Caller().Logger()

	return &LogUtil{
		logDir: logDir,
		logger: logger,
	}
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

func (lu *LogUtil) Error(message string) {
	lu.logger.Error().Msg(message)
}

func (lu *LogUtil) Errorf(format string, args ...interface{}) {
	lu.logger.Error().Msgf(format, args...)
}

// LogAndReturnError 记录错误日志并返回错误
func (lu *LogUtil) LogAndReturnError(message string) error {
	lu.Error(message)
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
