package webhook

import (
	"fmt"

	"servon/core/internal/libs/github/config"
	"servon/core/internal/libs/github/storage"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	config  *config.Manager
	storage *storage.Manager
}

func NewHandler(config *config.Manager) *Handler {
	return &Handler{
		config:  config,
		storage: storage.NewManager(),
	}
}

func (h *Handler) Handle(c *gin.Context) error {
	// 验证webhook签名
	signature := c.GetHeader("X-Hub-Signature-256")
	if signature == "" {
		return fmt.Errorf("missing signature")
	}

	// 读取请求体
	payload, err := c.GetRawData()
	if err != nil {
		return fmt.Errorf("failed to read payload")
	}

	// TODO: 实现签名验证

	event := c.GetHeader("X-GitHub-Event")
	eventID := c.GetHeader("X-GitHub-Delivery")

	// 处理不同类型的事件
	if err := h.handleEvent(event, payload); err != nil {
		return err
	}

	// 保存webhook数据
	return h.storage.SaveWebhookPayload(event, eventID, payload)
}

func (h *Handler) handleEvent(event string, payload []byte) error {
	switch event {
	case "installation", "installation_repositories":
		return h.handleInstallationEvent(payload)
	case "push":
		return h.handlePushEvent(payload)
	case "pull_request":
		return h.handlePullRequestEvent(payload)
	case "check_suite":
		return h.handleCheckSuiteEvent(payload)
	default:
		// 记录未处理的事件类型
		return nil
	}
}

// 实现各种事件处理方法...
func (h *Handler) handleInstallationEvent(payload []byte) error {
	// TODO: 实现安装事件处理逻辑
	return nil
}

func (h *Handler) handlePushEvent(payload []byte) error {
	// TODO: 实现推送事件处理逻辑
	return nil
}

func (h *Handler) handlePullRequestEvent(payload []byte) error {
	// TODO: 实现PR事件处理逻辑
	return nil
}

func (h *Handler) handleCheckSuiteEvent(payload []byte) error {
	// TODO: 实现检查套件事件处理逻辑
	return nil
}
