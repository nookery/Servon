package libs

var DefaultPrinter = NewPrinter()
var DefaultCommandManager = NewCommandManager()
var DefaultVersionManager = NewVersionManager()
var DefaultSystemResourcesManager = NewSystemResourcesManager()
var DefaultPortManager = NewPortManager()
var DefaultBasicInfoManager = newBasicInfoManager()
var DefaultProcessManager = NewProcessManager()
var DefaultFilesManager = NewFilesManager()
var DefaultOSInfoManager = NewOSInfoManager()
var DefaultNetworkManager = NewNetworkManager()
var DefaultServiceManager = NewServiceManager()
var DefaultAptManager = NewAptManager()
var DefaultDpkg = NewDpkg()
var DefaultCronManager = NewCronManager()
var DefaultSoftManager = newSoftManager()
var DefaultDeployManager = NewDeployManager()
var DefaultDataManager = NewDataManager()
var DefaultUserManager = NewUserManager()
var DefaultDownloadManager = NewDownloadManager()
var DefaultGitManager = NewGitManager()
var DefaultTaskManager = NewTaskManager()

// 以下为快捷函数

var GetOSInfo = DefaultOSInfoManager.GetOSInfo
var GetOSType = DefaultOSInfoManager.GetOSType

var NewCommand = DefaultCommandManager.NewCommand

var PrintInfo = DefaultPrinter.PrintInfo
var PrintLn = DefaultPrinter.PrintLn
var PrintInfof = DefaultPrinter.PrintInfof
var PrintInfofWithoutLocation = DefaultPrinter.PrintInfofWithoutLocation
var PrintInfofWithoutLocationAndEmoji = DefaultPrinter.PrintInfofWithoutLocationAndEmoji
var PrintTitle = DefaultPrinter.PrintTitle
var PrintSuccess = DefaultPrinter.PrintSuccess
var PrintSuccessf = DefaultPrinter.PrintSuccessf
var PrintWarn = DefaultPrinter.PrintWarn
var PrintWarnf = DefaultPrinter.PrintWarnf
var PrintError = DefaultPrinter.PrintError
var PrintErrorf = DefaultPrinter.PrintErrorf
var PrintErrorMessage = DefaultPrinter.PrintErrorMessage
var PrintList = DefaultPrinter.PrintList
var PrintCommand = DefaultPrinter.PrintCommand
var PrintCommandOutput = DefaultPrinter.PrintCommandOutput
var RunShell = DefaultPrinter.RunShell
var RunShellWithOutput = DefaultPrinter.RunShellWithOutput
