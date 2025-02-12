package libs

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
)

type Printer struct {
	Color   *color.Color
	LogChan chan string // 添加日志通道
	enabled bool        // 是否启用通道
}

func NewPrinter() *Printer {
	p := &Printer{
		Color:   color.New(color.FgCyan),
		LogChan: make(chan string, 100), // 使用带缓冲的通道，避免阻塞
		enabled: true,
	}

	// go func() {
	// 	for msg := range p.LogChan {
	// 		fmt.Println("=====:" + msg)
	// 	}
	// }()

	return p
}

// EnableLogging 启用日志通道
func (p *Printer) EnableLogging(enable bool) {
	p.enabled = enable
}

// Close 关闭日志通道
func (p *Printer) Close() {
	if p.LogChan != nil {
		close(p.LogChan)
	}
}

// sendToChannel 发送日志到通道
func (p *Printer) sendToChannel(message string) {
	if p.enabled && p.LogChan != nil {
		select {
		case p.LogChan <- message:
			// 成功发送
		default:
			// 通道已满，丢弃日志
		}
	}
}

// Print 打印信息
func (p *Printer) Print(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf(message)
	fmt.Println()
	p.sendToChannel(message)
}

// PrintCyan 打印青色信息
func (p *Printer) PrintCyan(format string, args ...interface{}) {
	p.Color.Printf(format, args...)
}

func (p *Printer) PrintGreen(format string, args ...interface{}) {
	p.Color.Printf(format, args...)
}

// PrintRed 打印红色信息
func (p *Printer) PrintRed(format string, args ...interface{}) {
	p.Color.Printf(format, args...)
}

// PrintWhite 打印白色信息
func (p *Printer) PrintWhite(format string, args ...interface{}) {
	p.Color.Printf(format, args...)
}

func (p *Printer) PrintYellow(format string, args ...interface{}) {
	p.Color.Printf(format, args...)
}

func (p *Printer) PrintAndReturnError(errMsg string) error {
	s := p.Color.Sprintf("❌ %s", errMsg)
	fmt.Println(s)
	return fmt.Errorf("%s", s)
}

// PrintInfo 打印信息
func (p *Printer) PrintInfo(format string, args ...interface{}) {
	message := fmt.Sprintf("🍋 "+format, args...)
	p.Color.Print(message)
	fmt.Println()
	p.sendToChannel(message)
}

// PrintInfof 打印信息
func (p *Printer) PrintInfof(format string, args ...interface{}) {
	p.Color.Printf("🍋 "+format, args...)
	fmt.Println()
	p.sendToChannel(fmt.Sprintf("🍋 "+format, args...))
}

// PrintInfofAndSend 打印信息并发送
func (p *Printer) PrintInfofAndSend(logChan chan<- string, format string, args ...interface{}) {
	message := fmt.Sprintf("🍋 "+format, args...)
	p.Color.Printf(message)
	fmt.Println()
	p.sendToChannel(message)
	logChan <- message
}

// PrintLn 打印换行
func (p *Printer) PrintLn() {
	p.Color.Println()
}

// Printf 打印格式化信息
func (p *Printer) Printf(format string, args ...interface{}) {
	p.Color.Printf(format, args...)
}

// PrintError 打印错误信息
func (p *Printer) PrintError(err error) {
	p.PrintErrorMessage(err.Error())
}

// PrintErrorf 打印错误信息
func (p *Printer) PrintErrorf(format string, args ...interface{}) {
	p.PrintErrorMessage(fmt.Sprintf(format, args...))
}

// PrintAndReturnErrorf 打印错误信息并返回错误
func (p *Printer) PrintAndReturnErrorf(format string, args ...interface{}) error {
	err := fmt.Errorf("%s", fmt.Sprintf(format, args...))
	p.PrintErrorMessage(err.Error())
	return err
}

// PrintErrorMessage 打印错误信息
func (p *Printer) PrintErrorMessage(message string) {
	_, thisFile, _, _ := runtime.Caller(0) // 获取当前文件路径
	_, file, line, _ := runtime.Caller(1)
	// 如果错误来自当前文件，则往上找一级调用者
	if file == thisFile {
		_, file, line, _ = runtime.Caller(2)
	}

	p.Color.Printf("❌ 错误: %s\n", message)
	p.Color.Printf("📃 位置: %s:%d\n", file, line)
	p.Color.Println()
	p.sendToChannel(fmt.Sprintf("❌ 错误: %s\n", message))
	p.sendToChannel(fmt.Sprintf("📃 位置: %s:%d\n", file, line))
}

// PrintList 打印列表
func (p *Printer) PrintList(list []string, title string) {
	p.Color.Println()
	p.Color.Println(title)
	if len(list) == 0 {
		color.New(color.FgYellow).Println("  暂无数据")
		fmt.Println()
		return
	}
	for _, item := range list {
		color.New(color.FgCyan).Printf("  ▶️  %s\n", item)
	}
	fmt.Println()
}

// PrintSuccess 打印成功信息
func (p *Printer) PrintSuccess(format string, args ...interface{}) {
	p.Color.Printf("✅ "+format, args...)
	p.Color.Println()
	p.sendToChannel(fmt.Sprintf("✅ "+format, args...))
}

// PrintWarn 打印警告信息
func (p *Printer) PrintWarn(format string, args ...interface{}) {
	p.Color.Printf("🚨 "+format, args...)
	p.Color.Println()
	p.sendToChannel(fmt.Sprintf("🚨 "+format, args...))
}

// PrintWarnf 打印警告信息
func (p *Printer) PrintWarnf(format string, args ...interface{}) {
	p.Color.Printf("🚨 "+format, args...)
	p.Color.Println()
	p.sendToChannel(fmt.Sprintf("🚨 "+format, args...))
}
