package events

import (
	"fmt"
	"sync"
)

// eventBus 私有结构体，实现IEventBus接口
type eventBus struct {
	subscribers     map[EventType][]Handler
	requestHandlers map[RequestType]RequestHandler
	mutex           sync.RWMutex
}

// 单例相关变量
var (
	instance IEventBus
	once     sync.Once
)

// GetEventBusInstance 获取EventBus的单例实例
// 这是获取EventBus实例的唯一公开方法
func GetEventBusInstance() IEventBus {
	once.Do(func() {
		instance = &eventBus{
			subscribers:     make(map[EventType][]Handler),
			requestHandlers: make(map[RequestType]RequestHandler),
		}
	})

	return instance
}

// Subscribe 订阅特定类型的事件
func (eb *eventBus) Subscribe(eventType EventType, handler Handler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

// Unsubscribe 取消订阅特定类型的事件
func (eb *eventBus) Unsubscribe(eventType EventType, handler Handler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if handlers, exists := eb.subscribers[eventType]; exists {
		for i, h := range handlers {
			if fmt.Sprintf("%v", h) == fmt.Sprintf("%v", handler) {
				eb.subscribers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

// Publish 发布事件
func (eb *eventBus) Publish(event Event) error {
	eb.mutex.RLock()
	handlers, exists := eb.subscribers[event.Type]
	eb.mutex.RUnlock()

	if !exists {
		return nil
	}

	// 异步调用所有处理器
	for _, handler := range handlers {
		go func(h Handler) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Recovered from panic in event handler: %v\n", r)
				}
			}()
			h(event)
		}(handler)
	}

	return nil
}

// RegisterRequestHandler 注册请求处理器
func (eb *eventBus) RegisterRequestHandler(requestType RequestType, handler RequestHandler) error {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if _, exists := eb.requestHandlers[requestType]; exists {
		return fmt.Errorf("handler already registered for request type: %s", requestType)
	}

	eb.requestHandlers[requestType] = handler
	return nil
}

// Request 发送同步请求并等待响应
func (eb *eventBus) Request(request Request) Response {
	eb.mutex.RLock()
	handler, exists := eb.requestHandlers[request.Type]
	eb.mutex.RUnlock()

	if !exists {
		return Response{
			Error: fmt.Sprintf("no handler registered for request type: %s", request.Type),
		}
	}

	return handler(request)
}
