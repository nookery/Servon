package libs

import (
	"fmt"

	"github.com/fatih/color"
)

// PrintCyan 打印青色信息
func PrintCyan(format string, args ...interface{}) {
	color.New(color.FgCyan).Printf(format, args...)
}

func PrintGreen(format string, args ...interface{}) {
	color.New(color.FgGreen).Printf(format, args...)
}

// PrintRed 打印红色信息
func PrintRed(format string, args ...interface{}) {
	color.New(color.FgRed).Printf(format, args...)
}

// PrintWhite 打印白色信息
func PrintWhite(format string, args ...interface{}) {
	color.New(color.FgWhite).Printf(format, args...)
}

func PrintYellow(format string, args ...interface{}) {
	color.New(color.FgYellow).Printf(format, args...)
}

func PrintAndReturnError(errMsg string) error {
	s := color.New(color.FgHiRed).Sprintf("❌ %s", errMsg)
	fmt.Println(s)
	return fmt.Errorf("%s", s)
}

// PrintInfo 打印信息
func PrintInfo(format string, args ...interface{}) {
	color.New(color.FgHiCyan).Printf("🍋 "+format, args...)
	fmt.Println()
}

// PrintLn 打印换行
func PrintLn() {
	fmt.Println()
}

// PrintError 打印错误信息
func PrintError(err error) {
	fmt.Println()
	color.New(color.FgHiRed).Printf("❌ 错误: %s\n", err.Error())
	fmt.Println()
}

// PrintErrorMessage 打印错误信息
func PrintErrorMessage(message string) {
	fmt.Println()
	color.New(color.FgHiRed).Printf("❌ 错误: %s\n", message)
	fmt.Println()
}

// PrintList 打印列表
func PrintList(list []string, title string) {
	fmt.Println()
	color.New(color.FgHiCyan).Println(title)
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
func PrintSuccess(format string, args ...interface{}) {
	color.New(color.FgHiGreen).Printf(format, args...)
	fmt.Println()
}
