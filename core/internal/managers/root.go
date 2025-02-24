package managers

import (
	"servon/core/internal/utils"
)

var logger = utils.DefaultLogUtil
var shell = utils.DefaultShellUtil

var PrintInfo = logger.Infof
var PrintErrorf = logger.Errorf
var PrintInfof = logger.Infof
var PrintSuccessf = logger.Successf
var PrintSuccess = logger.Success
var PrintError = logger.Error
var PrintList = logger.List
var PrintTitle = logger.Title
var PrintListWithTitle = logger.ListWithTitle
var PrintCommandOutput = shell.ExecuteWithOutput
var PrintAndReturnError = logger.LogAndReturnError
var PrintAndReturnErrorf = logger.LogAndReturnErrorf
var PrintAlert = logger.Alert
var PrintLn = logger.EmptyLine
var PrintErrorMessage = logger.ErrorMessage

var RunShell = shell.Execute
var RunShellWithSudo = shell.ExecuteWithSudo
var RunShellWithOutput = shell.RunShellWithOutput
