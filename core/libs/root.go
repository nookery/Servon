package libs

var DefaultPrinter = NewPrinter()
var DefaultLogManager = NewLogManager()
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

// 以下为快捷函数

var NewCommand = DefaultCommandManager.NewCommand
var StreamCommand = DefaultCommandManager.StreamCommand
var RunShell = DefaultShellManager.RunShell
var PrintInfo = DefaultPrinter.PrintInfo
var PrintInfof = DefaultPrinter.PrintInfof
var PrintSuccess = DefaultPrinter.PrintSuccess
var PrintError = DefaultPrinter.PrintError
var PrintErrorf = DefaultPrinter.PrintErrorf
var PrintList = DefaultPrinter.PrintList
