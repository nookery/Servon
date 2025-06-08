package commands

import (
	"servon/components/command_util"
	logger1 "servon/components/logger"
	"servon/components/shell_util"
	"servon/components/string_util"
)

var logger = logger1.DefaultLogUtil
var shell = shell_util.DefaultShellUtil
var stringUtil = string_util.DefaultStringUtil

var NewCommand = command_util.NewCommand
var RunShell = shell.RunShell
var PrintKeyValue = logger.PrintKeyValue
var PrintKeyValues = logger.PrintKeyValues
var PrintListWithTitle = logger.ListWithTitle
var PrintInfo = logger.Info
var PrintInfof = logger.Infof
var PrintCommandOutput = shell.RunShellWithOutput
var PrintSuccessf = logger.Successf
var PrintSuccess = logger.Success
var PrintError = logger.Error
var PrintErrorf = logger.Errorf
var PrintList = logger.List
var PrintTitle = logger.Title

type CommandOptions = command_util.CommandOptions
