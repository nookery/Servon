package managers

import (
	"servon/core/internal/utils"
)

var printer = utils.DefaultPrinter

var PrintInfo = printer.PrintInfo
var PrintErrorf = printer.PrintErrorf
var PrintInfof = printer.PrintInfof
var PrintSuccessf = printer.PrintSuccessf
var PrintSuccess = printer.PrintSuccess
var PrintError = printer.PrintError
var PrintList = printer.PrintList
var PrintTitle = printer.PrintTitle
var PrintCommandOutput = printer.PrintCommandOutput
var PrintAndReturnError = printer.PrintAndReturnError
var PrintAlert = printer.PrintAlert
var PrintLn = printer.PrintLn
var PrintErrorMessage = printer.PrintErrorMessage
var RunShell = printer.RunShell
var RunShellWithSudo = printer.RunShellWithSudo
var RunShellWithOutput = printer.RunShellWithOutput
