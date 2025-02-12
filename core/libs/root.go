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
var DefaultShellManager = NewShellManager()
var DefaultUserManager = NewUserManager()
var DefaultEnvManager = NewEnvManager()

// 以下为快捷函数

var GetOSInfo = DefaultOSInfoManager.GetOSInfo
var GetOSType = DefaultOSInfoManager.GetOSType

var NewCommand = DefaultCommandManager.NewCommand

var Print = DefaultPrinter.Print
var PrintInfo = DefaultPrinter.PrintInfo
var PrintInfof = DefaultPrinter.PrintInfof
var PrintSuccess = DefaultPrinter.PrintSuccess
var PrintError = DefaultPrinter.PrintError
var PrintErrorf = DefaultPrinter.PrintErrorf
var PrintList = DefaultPrinter.PrintList

var RunShell = DefaultShellManager.RunShell
var RunShellWithOutput = DefaultShellManager.RunShellWithOutput
var StreamCommand = DefaultCommandManager.StreamCommand
