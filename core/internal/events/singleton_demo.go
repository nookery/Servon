package events

import (
	"fmt"
	"os"
	"path/filepath"
)

// DemonstrateSingletonSafety 演示单例模式的安全性
// 这个函数展示了现在无法直接实例化eventBus结构体
func DemonstrateSingletonSafety() {
	fmt.Println("=== EventBus 单例模式安全性演示 ===")

	// ❌ 以下代码将无法编译，因为eventBus是私有结构体
	// var eb eventBus  // 编译错误：cannot refer to unexported name events.eventBus
	// eb := &eventBus{} // 编译错误：cannot refer to unexported name events.eventBus

	// ❌ 直接使用NewEventBus会返回错误
	if _, err := NewEventBus("/tmp/test"); err != nil {
		fmt.Printf("✓ NewEventBus 正确阻止了直接实例化: %v\n", err)
	}

	// ✅ 只能通过GetEventBusInstance获取实例
	tempDir := filepath.Join(os.TempDir(), "demo_events")
	defer os.RemoveAll(tempDir)

	instance1, err := GetEventBusInstance(tempDir)
	if err != nil {
		fmt.Printf("❌ 获取实例失败: %v\n", err)
		return
	}

	instance2, err := GetEventBusInstance(tempDir)
	if err != nil {
		fmt.Printf("❌ 获取实例失败: %v\n", err)
		return
	}

	if instance1 == instance2 {
		fmt.Println("✓ 两次调用GetEventBusInstance返回相同实例")
	} else {
		fmt.Println("❌ 单例模式失败：返回了不同的实例")
	}

	// ✅ 可以正常使用接口方法
	instance1.Subscribe(GitPush, func(e Event) {
		fmt.Printf("✓ 成功订阅事件: %s\n", e.Type)
	})

	err = instance1.Publish(Event{
		Type: GitPush,
		Data: map[string]interface{}{"repo": "demo-repo"},
	})
	if err != nil {
		fmt.Printf("❌ 发布事件失败: %v\n", err)
	} else {
		fmt.Println("✓ 成功发布事件")
	}

	fmt.Println("=== 演示完成 ===")
}