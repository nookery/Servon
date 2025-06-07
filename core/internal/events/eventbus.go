package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// EventBus 实现了一个事件总线，用于在系统中进行事件的发布-订阅和请求-响应通信。
// 它提供以下主要功能：
//   - 事件的发布与订阅（Publish/Subscribe）模式
//   - 同步请求处理（Request/Response）模式
//   - 自动的事件日志记录和历史查询
//
// 使用示例：
//
//	// 获取 EventBus 单例实例
//	eventBus, err := GetEventBusInstance("/path/to/logs")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// 订阅事件
//	eventBus.Subscribe(EventTypeUserCreated, func(e Event) {
//	    // 处理用户创建事件
//	})
//
//	// 发布事件
//	err = eventBus.Publish(Event{
//	    Type: EventTypeUserCreated,
//	    Data: map[string]interface{}{"userId": "123"},
//	})
//
//	// 注册请求处理器
//	eventBus.RegisterRequestHandler(RequestTypeGetUser, func(req Request) Response {
//	    // 处理获取用户请求
//	    return Response{Data: user}
//	})
//
//	// 发送请求
//	response := eventBus.Request(Request{
//	    Type: RequestTypeGetUser,
//	    Data: map[string]interface{}{"userId": "123"},
//	})
// IEventBus 定义事件总线接口
type IEventBus interface {
	Subscribe(eventType EventType, handler Handler)
	Unsubscribe(eventType EventType, handler Handler)
	Publish(event Event) error
	GetEventHistory(date time.Time) ([]Event, error)
	RegisterRequestHandler(requestType RequestType, handler RequestHandler) error
	Request(request Request) Response
}

// eventBus 私有结构体，实现IEventBus接口
type eventBus struct {
	subscribers     map[EventType][]Handler
	requestHandlers map[RequestType]RequestHandler
	mutex           sync.RWMutex
	logDir          string
}

// 单例相关变量
var (
	instance IEventBus
	once     sync.Once
)

// GetEventBusInstance 获取EventBus的单例实例
// 这是获取EventBus实例的唯一公开方法
func GetEventBusInstance(logDir string) (IEventBus, error) {
	var err error
	once.Do(func() {
		// 确保日志目录存在
		if mkdirErr := os.MkdirAll(logDir, 0755); mkdirErr != nil {
			err = fmt.Errorf("failed to create event log directory: %w", mkdirErr)
			return
		}

		instance = &eventBus{
			subscribers:     make(map[EventType][]Handler),
			requestHandlers: make(map[RequestType]RequestHandler),
			logDir:          logDir,
		}
	})

	if err != nil {
		return nil, err
	}

	// 如果实例已存在但logDir不同，更新logDir
	if eb, ok := instance.(*eventBus); ok && eb.logDir != logDir {
		eb.mutex.Lock()
		eb.logDir = logDir
		eb.mutex.Unlock()
		// 确保新的日志目录存在
		if mkdirErr := os.MkdirAll(logDir, 0755); mkdirErr != nil {
			return nil, fmt.Errorf("failed to create event log directory: %w", mkdirErr)
		}
	}

	return instance, nil
}

// NewEventBus 已弃用的构造函数，禁止外部直接创建实例
// 已弃用：请使用 GetEventBusInstance 方法获取单例实例
func NewEventBus(logDir string) (IEventBus, error) {
	return nil, fmt.Errorf("direct instantiation is not allowed, use GetEventBusInstance instead")
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
	// 记录事件到日志文件
	if err := eb.logEvent(event); err != nil {
		return fmt.Errorf("failed to log event: %w", err)
	}

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

// logEvent 将事件记录到日志文件
func (eb *eventBus) logEvent(event Event) error {
	// 创建日志文件名（按日期分割）
	logFile := filepath.Join(eb.logDir, fmt.Sprintf("events_%s.log", time.Now().Format("2006-01-02")))

	// 创建日志条目
	logEntry := struct {
		Timestamp time.Time `json:"timestamp"`
		Event     Event     `json:"event"`
	}{
		Timestamp: time.Now(),
		Event:     event,
	}

	// 序列化日志条目
	data, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// 追加到日志文件
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	return nil
}

// GetEventHistory 获取指定日期的事件历史
func (eb *eventBus) GetEventHistory(date time.Time) ([]Event, error) {
	logFile := filepath.Join(eb.logDir, fmt.Sprintf("events_%s.log", date.Format("2006-01-02")))

	data, err := os.ReadFile(logFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Event{}, nil
		}
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	var events []Event
	for _, line := range bytes.Split(data, []byte("\n")) {
		if len(line) == 0 {
			continue
		}

		var logEntry struct {
			Timestamp time.Time `json:"timestamp"`
			Event     Event     `json:"event"`
		}

		if err := json.Unmarshal(line, &logEntry); err != nil {
			return nil, fmt.Errorf("failed to unmarshal log entry: %w", err)
		}

		events = append(events, logEntry.Event)
	}

	return events, nil
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
