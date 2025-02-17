package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"servon/core/internal/libs/github/models"

	"github.com/gin-gonic/gin"
)

const baseURL = "http://43.142.208.212:9754" // TODO: 从环境变量或配置文件读取

type Manager struct {
	config *models.GitHubConfig
}

func NewGitHubConfig() *Manager {
	return &Manager{
		config: &models.GitHubConfig{
			Installations: make(map[int64]*models.Installation),
		},
	}
}

func (m *Manager) HandleSetup(c *gin.Context) (string, error) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		return "", fmt.Errorf("invalid request: %v", err)
	}

	manifest := createManifest(req.Name, req.Description)
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return "", fmt.Errorf("failed to generate manifest: %v", err)
	}

	state, err := generateState()
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %v", err)
	}

	return generateHTML(state, string(manifestJSON)), nil
}

func (m *Manager) HandleCallback(c *gin.Context) (string, error) {
	code := c.Query("code")
	if code == "" {
		return "", fmt.Errorf("missing code parameter")
	}

	result, err := convertManifest(code)
	if err != nil {
		return "", err
	}

	// Update config
	m.config.GitHubAppID = result.ID
	m.config.GitHubAppPrivateKey = result.PEM

	if err := m.SaveConfig(); err != nil {
		return "", fmt.Errorf("failed to save config: %v", err)
	}

	return fmt.Sprintf("https://github.com/apps/%s/installations/new", result.Name), nil
}

func (m *Manager) SaveConfig() error {
	// TODO: Implement persistence logic
	return nil
}

// Private helper functions

func createManifest(name, description string) models.GitHubManifest {
	manifest := models.GitHubManifest{
		Name:        name,
		URL:         baseURL,
		Description: description,
		Public:      true,
		HookAttributes: struct {
			URL    string `json:"url"`
			Active bool   `json:"active"`
		}{
			URL:    fmt.Sprintf("%s/web_api/github/webhook", baseURL),
			Active: true,
		},
		RedirectURL:  fmt.Sprintf("%s/web_api/github/callback", baseURL),
		CallbackURLs: []string{fmt.Sprintf("%s/web_api/github/callback", baseURL)},
		DefaultPermissions: map[string]string{
			"issues": "write",
			"checks": "write",
		},
		DefaultEvents: []string{
			"issues",
			"issue_comment",
			"check_suite",
			"check_run",
		},
	}
	return manifest
}

func generateState() (string, error) {
	state := make([]byte, 16)
	if _, err := rand.Read(state); err != nil {
		return "", err
	}
	return hex.EncodeToString(state), nil
}

func generateHTML(state, manifestJSON string) string {
	return fmt.Sprintf(`
		<form id="github-form" action="https://github.com/settings/apps/new?state=%s" method="post">
			<input type="hidden" name="manifest" value='%s'>
		</form>
		<script>document.getElementById("github-form").submit();</script>
	`, state, manifestJSON)
}

func convertManifest(code string) (*struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	WebhookURL string `json:"webhook_url"`
	PEM        string `json:"pem"`
}, error) {
	resp, err := http.Post(
		fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions", code),
		"application/json",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub App: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		ID         int64  `json:"id"`
		Name       string `json:"name"`
		WebhookURL string `json:"webhook_url"`
		PEM        string `json:"pem"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub response: %v", err)
	}

	return &result, nil
}
