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

// EventBus 实现事件总线
type EventBus struct {
	subscribers     map[EventType][]Handler
	requestHandlers map[RequestType]RequestHandler
	mutex           sync.RWMutex
	logDir          string
}

// NewEventBus 创建新的事件总线实例
func NewEventBus(logDir string) (*EventBus, error) {
	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create event log directory: %w", err)
	}

	return &EventBus{
		subscribers:     make(map[EventType][]Handler),
		requestHandlers: make(map[RequestType]RequestHandler),
		logDir:          logDir,
	}, nil
}

// Subscribe 订阅特定类型的事件
func (eb *EventBus) Subscribe(eventType EventType, handler Handler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

// Unsubscribe 取消订阅特定类型的事件
func (eb *EventBus) Unsubscribe(eventType EventType, handler Handler) {
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
func (eb *EventBus) Publish(event Event) error {
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
func (eb *EventBus) logEvent(event Event) error {
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
func (eb *EventBus) GetEventHistory(date time.Time) ([]Event, error) {
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
func (eb *EventBus) RegisterRequestHandler(requestType RequestType, handler RequestHandler) error {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if _, exists := eb.requestHandlers[requestType]; exists {
		return fmt.Errorf("handler already registered for request type: %s", requestType)
	}

	eb.requestHandlers[requestType] = handler
	return nil
}

// Request 发送同步请求并等待响应
func (eb *EventBus) Request(request Request) Response {
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
