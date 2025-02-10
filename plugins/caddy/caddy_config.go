package caddy

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// 定义模板文件的路径常量
const (
	caddyBaseTemplate = "templates/caddy_base.conf"
	caddySiteTemplate = "templates/caddy_site.conf.tmpl"
)

//go:embed templates/caddy_site.conf.tmpl templates/caddy_base.conf
var templateFS embed.FS

// CaddyConfig 处理所有与 Caddy 相关的配置操作
type CaddyConfig struct {
	BaseDir string
}

// GetConfigDir 返回 Caddy 的配置目录
func (cc *CaddyConfig) GetConfigDir() string {
	return cc.BaseDir
}

// GetCaddyfilePath 返回主 Caddyfile 的路径
func (cc *CaddyConfig) GetCaddyfilePath() string {
	return filepath.Join(cc.GetConfigDir(), "Caddyfile")
}

// GetProjectConfigPath 返回特定项目的配置文件路径
func (cc *CaddyConfig) GetProjectConfigPath(projectName string) string {
	return filepath.Join(cc.GetConfigDir(), fmt.Sprintf("%s.conf", projectName))
}

// EnsureConfigDir 确保配置文件的存储目录存在
func (cc *CaddyConfig) EnsureConfigDir() error {
	return os.MkdirAll(cc.GetConfigDir(), 0755)
}

// EnsureCaddyfile 确保主 Caddyfile 存在，如果需要则从模板创建
func (cc *CaddyConfig) EnsureCaddyfile() error {
	caddyfile := cc.GetCaddyfilePath()
	if _, err := os.Stat(caddyfile); os.IsNotExist(err) {
		templateCaddyfile, err := templateFS.ReadFile(caddyBaseTemplate)
		if err != nil {
			return fmt.Errorf("failed to read Caddyfile template: %v", err)
		}

		if err := os.WriteFile(caddyfile, templateCaddyfile, 0644); err != nil {
			return fmt.Errorf("failed to create Caddyfile: %v", err)
		}
	}
	return nil
}

// UpdateProjectConfig 更新特定项目的配置
func (cc *CaddyConfig) UpdateProjectConfig(project *Project) error {
	// 确保配置目录存在
	if err := cc.EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// 读取站点配置模板
	templateContent, err := templateFS.ReadFile(caddySiteTemplate)
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

// WriteConfig 将配置内容写入指定路径
// path: 配置文件路径
// content: 配置文件内容
func (cc *CaddyConfig) WriteConfig(path string, content string) error {
	// 确保配置目录存在
	if err := cc.EnsureConfigDir(); err != nil {
		return fmt.Errorf("确保配置目录存在失败: %v", err)
	}

	// 写入配置文件
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}

// RemoveConfig 删除指定路径的配置文件
// path: 要删除的配置文件路径
func (cc *CaddyConfig) RemoveConfig(path string) error {
	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // 文件不存在则直接返回
	}

	// 删除配置文件
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("删除配置文件失败: %v", err)
	}

	return nil
}
