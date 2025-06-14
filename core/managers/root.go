package managers

import (
	"servon/components/shell_util"
	"servon/components/utils"
)

var logger = utils.DefaultLogUtil
var shell = shell_util.DefaultShellUtil

var PrintInfo = logger.Infof
var PrintErrorf = logger.Errorf
var PrintInfof = logger.Infof
var PrintSuccessf = logger.Successf
var PrintSuccess = logger.Success
var PrintError = logger.Error
var PrintList = logger.List
var PrintTitle = logger.Title
var PrintListWithTitle = logger.ListWithTitle
var PrintCommandOutput = shell.RunShellWithOutput
var PrintAndReturnError = logger.LogAndReturnError
var PrintAndReturnErrorf = logger.LogAndReturnErrorf
var PrintAlert = logger.Alert
var PrintLn = logger.EmptyLine
var PrintErrorMessage = logger.ErrorMessage

var RunShell = shell.RunShell
var RunShellWithSudo = shell.RunShellWithSudo
var RunShellWithOutput = shell.RunShellWithOutput
