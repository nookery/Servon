// Package components 提供相对独立的功能组件
//
// 这个包的设计理念是将系统中相对独立、可复用的功能组件集中管理。
// 每个子包都实现特定的功能领域，具有以下特点：
//
// 1. **独立性**: 各组件之间耦合度低，可以独立开发和测试
// 2. **可复用性**: 组件设计通用，可在不同场景下复用
// 3. **统一入口**: 通过root.go文件统一对外提供组件实例
//
// 使用示例：
//
//	// 使用事件总线
//	components.EventBus.Subscribe(events.EventTypeUserCreated, handler)
//	components.EventBus.Publish(event)
//
// 添加新组件的指导原则：
//  1. 确保组件功能相对独立
//  2. 提供清晰的接口定义
//  3. 在root.go中暴露组件实例
//  4. 编写完整的测试用例
package components

import "servon/components/events"

// EventBus 全局事件总线实例
// 提供系统级的事件发布-订阅和请求-响应功能
var EventBus = events.GetEventBusInstance()
