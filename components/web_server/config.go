package web_server

const PID_FILE = "servon.pid"
const LOG_FILE = "servon.log"
const DEFAULT_PORT = 9876

// WebServerConfig 服务器配置
type WebServerConfig struct {
	Host    string
	Port    int
	Verbose bool
	Logger  Logger // 自定义日志器，如果为 nil 则使用默认日志器
}
