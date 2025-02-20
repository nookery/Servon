package commands

import (
	"servon/core/internal/libs/utils"
)

var printer = utils.DefaultPrinter
var stringUtil = utils.DefaultStringUtil
var NewCommand = utils.NewCommand
var RunShell = printer.RunShell
var PrintKeyValue = printer.PrintKeyValue
var PrintKeyValues = printer.PrintKeyValues
var PrintInfo = printer.PrintInfo
var PrintErrorf = printer.PrintErrorf
var PrintInfof = printer.PrintInfof
var PrintCommandOutput = printer.PrintCommandOutput
var PrintSuccessf = printer.PrintSuccessf
var PrintSuccess = printer.PrintSuccess
var PrintError = printer.PrintError
var PrintList = printer.PrintList
var PrintTitle = printer.PrintTitle

type CommandOptions = utils.CommandOptions
