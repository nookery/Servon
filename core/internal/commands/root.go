package commands

import (
	"servon/core/internal/utils"
)

var logger = utils.DefaultLogUtil
var shell = utils.DefaultShellUtil
var stringUtil = utils.DefaultStringUtil

var NewCommand = utils.NewCommand
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

type CommandOptions = utils.CommandOptions
