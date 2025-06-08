package events

// IEventBus 定义事件总线接口
// 
// 使用示例：
//
//	// 获取 EventBus 单例实例
//	eventBus := GetEventBusInstance()
//
//	// 订阅事件
//	eventBus.Subscribe(EventTypeUserCreated, func(e Event) {
//	    // 处理用户创建事件
//	})
//
//	// 发布事件
//	eventBus.Publish(Event{
//	    Type: EventTypeUserCreated,
//	    Data: map[string]interface{}{"userId": "123"},
//	})
type IEventBus interface {
	Subscribe(eventType EventType, handler Handler)
	Unsubscribe(eventType EventType, handler Handler)
	Publish(event Event) error
	RegisterRequestHandler(requestType RequestType, handler RequestHandler) error
	Request(request Request) Response
}
