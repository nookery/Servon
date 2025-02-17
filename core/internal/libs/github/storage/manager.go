package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"servon/core/internal/libs/github/models"
)

const dataDir = "/data/github"

type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) SaveWebhookPayload(eventType, eventID string, payload []byte) error {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	filename := fmt.Sprintf("%s/%d_%s_%s.json",
		dataDir,
		time.Now().Unix(),
		eventID,
		eventType,
	)

	if err := os.WriteFile(filename, payload, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

func (m *Manager) GetWebhooks() ([]models.WebhookPayload, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create webhooks directory: %v", err)
	}

	files, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read webhooks directory: %v", err)
	}

	var webhooks []models.WebhookPayload
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			webhook, err := m.readWebhookFile(file.Name())
			if err != nil {
				continue
			}
			webhooks = append(webhooks, webhook)
		}
	}

	return webhooks, nil
}

func (m *Manager) readWebhookFile(filename string) (models.WebhookPayload, error) {
	var webhook models.WebhookPayload

	parts := strings.Split(strings.TrimSuffix(filename, ".json"), "_")
	if len(parts) != 3 {
		return webhook, fmt.Errorf("invalid filename format")
	}

	data, err := os.ReadFile(filepath.Join(dataDir, filename))
	if err != nil {
		return webhook, err
	}

	var payload interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return webhook, err
	}

	webhook.Timestamp = parseTimestamp(parts[0])
	webhook.ID = parts[1]
	webhook.Type = parts[2]
	webhook.Payload = payload

	return webhook, nil
}

func parseTimestamp(ts string) int64 {
	timestamp, _ := time.Parse(time.RFC3339, ts)
	return timestamp.Unix()
}
