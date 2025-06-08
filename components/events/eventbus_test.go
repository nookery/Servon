package events

import (
	"testing"
	"time"
)

// TestSingletonPattern 测试单例模式
func TestSingletonPattern(t *testing.T) {
	// 获取第一个实例
	instance1 := GetEventBusInstance()

	// 获取第二个实例
	instance2 := GetEventBusInstance()

	// 验证两个实例是同一个对象
	if instance1 != instance2 {
		t.Error("Expected same instance, got different instances")
	}

	// 验证实例不为nil
	if instance1 == nil {
		t.Error("Instance should not be nil")
	}
}

// TestEventBusBasicFunctionality 测试EventBus基本功能
func TestEventBusBasicFunctionality(t *testing.T) {
	// 获取EventBus实例
	eventBus := GetEventBusInstance()

	// 测试事件订阅和发布
	eventReceived := false
	eventBus.Subscribe(GitPush, func(e Event) {
		eventReceived = true
	})

	// 发布事件
	err := eventBus.Publish(Event{
		Type: GitPush,
		Data: map[string]interface{}{"repo": "test-repo"},
	})
	if err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	// 等待事件处理（异步）
	time.Sleep(100 * time.Millisecond)

	if !eventReceived {
		t.Error("Event was not received")
	}

	// 测试请求处理
	err = eventBus.RegisterRequestHandler(SoftwareInfoRequest, func(req Request) Response {
		return Response{
			Data: map[string]interface{}{"software": "test"},
		}
	})
	if err != nil {
		t.Fatalf("Failed to register request handler: %v", err)
	}

	// 发送请求
	response := eventBus.Request(Request{
		Type: SoftwareInfoRequest,
		Data: map[string]interface{}{"name": "test-software"},
	})

	if response.Error != "" {
		t.Errorf("Unexpected error in response: %s", response.Error)
	}

	if response.Data == nil {
		t.Error("Expected response data, got nil")
	}
}
