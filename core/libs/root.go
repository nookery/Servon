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

var NewCommand = DefaultCommandManager.NewCommand
var StreamCommand = DefaultCommandManager.StreamCommand
