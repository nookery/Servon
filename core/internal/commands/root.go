package commands

import (
	"servon/core/internal/libs"
	"servon/core/internal/utils"
)

var DefaultPrinter = utils.DefaultPrinter
var NewCommand = utils.NewCommand
var RunShell = libs.RunShell
var PrintKeyValue = DefaultPrinter.PrintKeyValue
var PrintKeyValues = DefaultPrinter.PrintKeyValues
var PrintInfo = DefaultPrinter.PrintInfo
var PrintErrorf = DefaultPrinter.PrintErrorf
var PrintInfof = DefaultPrinter.PrintInfof
var PrintCommandOutput = DefaultPrinter.PrintCommandOutput
var PrintSuccessf = DefaultPrinter.PrintSuccessf
var PrintSuccess = DefaultPrinter.PrintSuccess
var PrintError = DefaultPrinter.PrintError
var PrintList = DefaultPrinter.PrintList
var PrintTitle = DefaultPrinter.PrintTitle

type CommandOptions = utils.CommandOptions
