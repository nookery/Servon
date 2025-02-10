package api

import (
	"fmt"
	"servon/core/libs"
)

type PrintApi struct{}

func NewPrintApi() PrintApi {
	return PrintApi{}
}

func (p *PrintApi) PrintLn() {
	libs.PrintLn()
}

func (p *PrintApi) PrintError(err error) {
	libs.PrintError(err)
}

func (p *PrintApi) PrintInfo(format string, args ...interface{}) {
	libs.PrintInfo(format, args...)
}

func (p *PrintApi) PrintCyan(format string, args ...interface{}) {
	libs.PrintCyan(format, args...)
}

func (p *PrintApi) PrintGreen(format string, args ...interface{}) {
	libs.PrintGreen(format, args...)
}

func (p *PrintApi) PrintRed(format string, args ...interface{}) {
	libs.PrintRed(format, args...)
}

func (p *PrintApi) PrintWhite(format string, args ...interface{}) {
	libs.PrintWhite(format, args...)
}

func (p *PrintApi) PrintYellow(format string, args ...interface{}) {
	libs.PrintYellow(format, args...)
}

func (p *PrintApi) PrintSuccess(format string, args ...interface{}) {
	libs.PrintSuccess(fmt.Sprintf("✅ %s", format), args...)
}
