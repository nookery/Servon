package libs

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
)

type Printer struct {
	Color *color.Color
}

func NewPrinter() *Printer {
	return &Printer{
		Color: color.New(color.FgCyan),
	}
}

// Print 打印信息
func (p *Printer) Print(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Println()
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
	p.Color.Printf("🍋 "+format, args...)
	fmt.Println()
}

// PrintInfof 打印信息
func (p *Printer) PrintInfof(format string, args ...interface{}) {
	p.Color.Printf("🍋 "+format, args...)
	fmt.Println()
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
}

// PrintWarn 打印警告信息
func (p *Printer) PrintWarn(format string, args ...interface{}) {
	p.Color.Printf("🚨 "+format, args...)
	p.Color.Println()
}

// PrintWarnf 打印警告信息
func (p *Printer) PrintWarnf(format string, args ...interface{}) {
	p.Color.Printf("🚨 "+format, args...)
	p.Color.Println()
}
