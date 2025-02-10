package caddy

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

type CaddyTemplate struct {
}

// RenderProxyConfig 渲染代理配置模板
func (tm *CaddyTemplate) RenderProxyConfig(domain string, target string) (string, error) {
	// 读取代理配置模板
	tmplContent, err := tm.ReadTemplate("templates/proxy.tmpl")
	if err != nil {
		return "", fmt.Errorf("读取代理配置模板失败: %v", err)
	}

	// 替换模板变量
	proxyConfig := strings.ReplaceAll(tmplContent, "{{domain}}", domain)
	proxyConfig = strings.ReplaceAll(proxyConfig, "{{target}}", target)

	return proxyConfig, nil
}

// ReadTemplate 从嵌入的模板文件系统读取指定的模板文件
func (tm *CaddyTemplate) ReadTemplate(path string) (string, error) {
	content, err := templatesFS.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("读取模板文件失败 %s: %v", path, err)
	}

	return string(content), nil
}
