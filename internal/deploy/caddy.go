package deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"servon/utils"
	"text/template"
)

const caddyConfigTemplate = `
{{ .Domain }} {
	{{ if eq .Type "static" }}
	root * {{ .OutputPath }}
	file_server
	{{ else }}
	reverse_proxy localhost:{{ .Port }}
	{{ end }}
}
`

// updateCaddyConfig 更新 Caddy 配置
func updateCaddyConfig(project *Project) error {
	utils.Info("更新 Caddy 配置: [%d] %s", project.ID, project.Domain)

	// 创建配置目录
	configDir := filepath.Join("data", "caddy", "conf.d")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		utils.Error("创建配置目录失败: %v", err)
		return fmt.Errorf("创建配置目录失败: %v", err)
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
		OutputPath: filepath.Join("data", "projects", fmt.Sprintf("%d", project.ID), project.OutputDir),
		Port:       project.Port,
	}

	// 解析并执行模板
	tmpl, err := template.New("caddy").Parse(caddyConfigTemplate)
	if err != nil {
		utils.Error("解析配置模板失败: %v", err)
		return fmt.Errorf("解析配置模板失败: %v", err)
	}

	// 创建配置文件
	configFile := filepath.Join(configDir, fmt.Sprintf("%d.conf", project.ID))
	f, err := os.Create(configFile)
	if err != nil {
		utils.Error("创建配置文件失败: %v", err)
		return fmt.Errorf("创建配置文件失败: %v", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		utils.Error("生成配置文件失败: %v", err)
		return fmt.Errorf("生成配置文件失败: %v", err)
	}

	// 重载 Caddy 配置
	if err := reloadCaddy(); err != nil {
		utils.Error("重载 Caddy 配置失败: %v", err)
		return fmt.Errorf("重载 Caddy 配置失败: %v", err)
	}

	utils.Info("Caddy 配置更新成功: %s", project.Domain)
	return nil
}

// reloadCaddy 重载 Caddy 配置
func reloadCaddy() error {
	// TODO: 实现 Caddy 配置重载逻辑
	return nil
}

// GetDomains 获取所有域名配置
func GetDomains() ([]string, error) {
	projectsMu.RLock()
	defer projectsMu.RUnlock()

	domains := make([]string, 0, len(projects))
	for _, p := range projects {
		domains = append(domains, p.Domain)
	}
	return domains, nil
}
