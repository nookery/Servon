package managers

import (
	"servon/core/internal/utils"
)

var DefaultPrinter = utils.DefaultPrinter

var PrintInfo = DefaultPrinter.PrintInfo
var PrintErrorf = DefaultPrinter.PrintErrorf
var PrintInfof = DefaultPrinter.PrintInfof
var PrintSuccessf = DefaultPrinter.PrintSuccessf
var PrintSuccess = DefaultPrinter.PrintSuccess
var PrintError = DefaultPrinter.PrintError
var PrintList = DefaultPrinter.PrintList
var PrintTitle = DefaultPrinter.PrintTitle
var PrintCommandOutput = DefaultPrinter.PrintCommandOutput
var PrintAndReturnError = DefaultPrinter.PrintAndReturnError
var PrintAlert = DefaultPrinter.PrintAlert
var PrintLn = DefaultPrinter.PrintLn
var PrintErrorMessage = DefaultPrinter.PrintErrorMessage
var RunShell = DefaultPrinter.RunShell
var RunShellWithSudo = DefaultPrinter.RunShellWithSudo
var RunShellWithOutput = DefaultPrinter.RunShellWithOutput
