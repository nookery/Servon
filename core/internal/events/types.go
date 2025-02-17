package events

// EventType 定义事件类型
type EventType string

// Event 定义事件结构
type Event struct {
	Type EventType   `json:"type"`
	Data interface{} `json:"data"`
}

// RequestType 定义请求类型
type RequestType string

// Request 定义请求结构
type Request struct {
	Type RequestType `json:"type"`
	Data interface{} `json:"data"`
}

// Response 定义响应结构
type Response struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error,omitempty"`
}

// RequestHandler 定义请求处理函数类型
type RequestHandler func(Request) Response

// 系统事件类型定义
const (
	// Git相关事件
	GitPush  EventType = "git:push"
	GitClone EventType = "git:clone"
	GitPull  EventType = "git:pull"

	// 部署相关事件
	DeployStart    EventType = "deploy:start"
	DeployComplete EventType = "deploy:complete"
	DeployFailed   EventType = "deploy:failed"

	// 服务相关事件
	ServiceStart   EventType = "service:start"
	ServiceStop    EventType = "service:stop"
	ServiceRestart EventType = "service:restart"

	// 软件相关事件
	SoftwareInstall   EventType = "software:install"
	SoftwareUninstall EventType = "software:uninstall"
	SoftwareUpgrade   EventType = "software:upgrade"
)

// 系统请求类型定义
const (
	// 软件相关请求
	SoftwareInfoRequest RequestType = "software:info"

	// 服务相关请求
	ServiceStatusRequest RequestType = "service:status"

	// 部署相关请求
	DeployStatusRequest RequestType = "deploy:status"

	// Git相关请求
	GitRepoInfoRequest RequestType = "git:repo:info"
)

// Handler 定义事件处理函数类型
type Handler func(Event)
