package libs

import (
	"fmt"

	"github.com/fatih/color"
)

type Printer struct {
	Color *color.Color
}

func NewPrint() *Printer {
	return &Printer{
		Color: color.New(color.FgCyan),
	}
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

// PrintLn 打印换行
func (p *Printer) PrintLn() {
	p.Color.Println()
}

// PrintError 打印错误信息
func (p *Printer) PrintError(err error) {
	p.Color.Println()
	p.Color.Printf("❌ 错误: %s\n", err.Error())
	p.Color.Println()
}

// PrintErrorMessage 打印错误信息
func (p *Printer) PrintErrorMessage(message string) {
	p.Color.Println()
	p.Color.Printf("❌ 错误: %s\n", message)
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
	p.Color.Printf(format, args...)
	p.Color.Println()
}
