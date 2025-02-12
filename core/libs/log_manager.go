package libs

import (
	"fmt"

	"github.com/fatih/color"
)

type LogManager struct {
}

func NewLogManager() *LogManager {
	return &LogManager{}
}

func (l *LogManager) Debug(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l *LogManager) Error(args ...interface{}) {
	fmt.Println(append([]interface{}{color.RedString("❌")}, args...)...)
}

func (l *LogManager) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Printf(format, args...)
}

func (l *LogManager) ErrorChan(ch chan<- string, format string, args ...interface{}) {
	ch <- fmt.Sprintf(format, args...)
}

// Info 打印信息
func (l *LogManager) Info(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Infoln 打印信息并换行
func (l *LogManager) Infoln(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// InfoWithSpace 打印信息并换行
func (l *LogManager) InfoWithSpace(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// InfoChan 打印信息到通道
func (l *LogManager) InfoChan(ch chan<- string, format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
	ch <- fmt.Sprintf(format, args...)
}
