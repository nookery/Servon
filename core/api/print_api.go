package api

import (
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
