package libs

import "servon/core/internal/utils"

var DefaultPrinter = utils.DefaultPrinter
var DefaultSystemResourcesManager = NewSystemResourcesManager()
var DefaultPortManager = NewPortManager()
var DefaultBasicInfoManager = newBasicInfoManager()
var DefaultOSInfoManager = NewOSInfoManager()
var DefaultNetworkManager = NewNetworkManager()
var DefaultDpkg = NewDpkg()
var DefaultUserManager = NewUserLib()

// 以下为快捷函数

var GetOSInfo = DefaultOSInfoManager.GetOSInfo
var GetOSType = DefaultOSInfoManager.GetOSType

var PrintInfo = DefaultPrinter.PrintInfo
var PrintLn = DefaultPrinter.PrintLn
var PrintInfof = DefaultPrinter.PrintInfof
var PrintInfofWithoutLocation = DefaultPrinter.PrintInfofWithoutLocation
var PrintInfofWithoutLocationAndEmoji = DefaultPrinter.PrintInfofWithoutLocationAndEmoji
var PrintTitle = DefaultPrinter.PrintTitle
var PrintSuccess = DefaultPrinter.PrintSuccess
var PrintSuccessf = DefaultPrinter.PrintSuccessf
var PrintWarn = DefaultPrinter.PrintWarn
var PrintAlert = DefaultPrinter.PrintAlert
var PrintWarnf = DefaultPrinter.PrintWarnf
var PrintError = DefaultPrinter.PrintError
var PrintErrorf = DefaultPrinter.PrintErrorf
var PrintErrorMessage = DefaultPrinter.PrintErrorMessage
var PrintList = DefaultPrinter.PrintList
var PrintCommand = DefaultPrinter.PrintCommand
var PrintCommandOutput = DefaultPrinter.PrintCommandOutput
var PrintAndReturnError = DefaultPrinter.PrintAndReturnError
var RunShell = DefaultPrinter.RunShell
var RunShellWithOutput = DefaultPrinter.RunShellWithOutput
