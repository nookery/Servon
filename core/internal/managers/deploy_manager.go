package managers

import (
	"log"

	"servon/core/internal/events"
)

type DeployManager struct {
	eventBus *events.EventBus
}

func NewDeployManager(eventBus *events.EventBus) *DeployManager {
	dm := &DeployManager{
		eventBus: eventBus,
	}

	// 订阅Git Push事件
	eventBus.Subscribe(events.GitPush, dm.handleGitPushEvent)

	return dm
}

// handleGitPushEvent 处理Git Push事件
func (m *DeployManager) handleGitPushEvent(event events.Event) {
	deployData, ok := event.Data.(map[string]interface{})
	if !ok {
		log.Printf("Invalid deploy data format")
		return
	}

	repo, ok := deployData["repository"].(string)
	if !ok {
		log.Printf("Repository information missing")
		return
	}

	// 执行部署操作
	if err := m.deployProject(repo); err != nil {
		log.Printf("Deploy failed: %v", err)

		// 发布部署失败事件
		m.eventBus.Publish(events.Event{
			Type: events.DeployFailed,
			Data: map[string]interface{}{
				"repository": repo,
				"error":      err.Error(),
			},
		})
		return
	}

	// 发布部署成功事件
	m.eventBus.Publish(events.Event{
		Type: events.DeployComplete,
		Data: map[string]interface{}{
			"repository": repo,
			"status":     "success",
		},
	})
}

// deployProject 执行实际的部署操作
func (m *DeployManager) deployProject(repo string) error {
	log.Printf("Deploying project from repository: %s", repo)

	// 发布部署开始事件
	m.eventBus.Publish(events.Event{
		Type: events.DeployStart,
		Data: map[string]interface{}{
			"repository": repo,
		},
	})

	// TODO: 实现实际的部署逻辑
	// 1. 拉取代码
	// 2. 构建项目
	// 3. 部署服务

	return nil
}
