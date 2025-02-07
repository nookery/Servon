package softwares

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// CaddyConfig handles all configuration related operations for Caddy
type CaddyConfig struct {
	baseDir string
}

// NewCaddyConfig creates a new CaddyConfig instance
func NewCaddyConfig() *CaddyConfig {
	return &CaddyConfig{
		baseDir: filepath.Join("data", "caddy"),
	}
}

// GetConfigDir returns the base configuration directory for Caddy
func (cc *CaddyConfig) GetConfigDir() string {
	return cc.baseDir
}

// GetCaddyfilePath returns the path to the main Caddyfile
func (cc *CaddyConfig) GetCaddyfilePath() string {
	return filepath.Join(cc.GetConfigDir(), "Caddyfile")
}

// GetProjectConfigPath returns the configuration file path for a specific project
func (cc *CaddyConfig) GetProjectConfigPath(projectName string) string {
	return filepath.Join(cc.GetConfigDir(), fmt.Sprintf("%s.conf", projectName))
}

// EnsureConfigDir ensures the Caddy configuration directory exists
func (cc *CaddyConfig) EnsureConfigDir() error {
	return os.MkdirAll(cc.GetConfigDir(), 0755)
}

// EnsureCaddyfile ensures the main Caddyfile exists, creating it from template if needed
func (cc *CaddyConfig) EnsureCaddyfile() error {
	caddyfile := cc.GetCaddyfilePath()
	if _, err := os.Stat(caddyfile); os.IsNotExist(err) {
		templateCaddyfile, err := os.ReadFile("internal/softwares/templates/caddy_base.conf")
		if err != nil {
			return fmt.Errorf("failed to read Caddyfile template: %v", err)
		}

		if err := os.WriteFile(caddyfile, templateCaddyfile, 0644); err != nil {
			return fmt.Errorf("failed to create Caddyfile: %v", err)
		}
	}
	return nil
}

// UpdateProjectConfig updates the configuration for a specific project
func (cc *CaddyConfig) UpdateProjectConfig(project *Project) error {
	// 确保配置目录存在
	if err := cc.EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// 读取站点配置模板
	templateContent, err := os.ReadFile("internal/softwares/templates/caddy_site.conf.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read site template: %v", err)
	}

	// 解析并执行模板
	tmpl, err := template.New("caddy").Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse config template: %v", err)
	}

	// 准备模板数据
	data := struct {
		Domain     string
		Type       string
		OutputPath string
		Port       int
	}{
		Domain:     project.Domain,
		Type:       project.Type,
		OutputPath: project.OutputDir,
		Port:       project.Port,
	}

	// 确保主 Caddyfile 存在
	if err := cc.EnsureCaddyfile(); err != nil {
		return fmt.Errorf("failed to ensure Caddyfile exists: %v", err)
	}

	// 创建项目配置文件
	configFile := cc.GetProjectConfigPath(project.Name)
	f, err := os.Create(configFile)
	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to generate config file: %v", err)
	}

	return nil
}
