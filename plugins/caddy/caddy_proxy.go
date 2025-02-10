package caddy

import (
	"fmt"
	"strings"
)

// AddProxy 添加反向代理配置
// domain: 要代理的域名
// target: 目标地址（例如：127.0.0.1:8888）
func (c *Caddy) AddProxy(domain string, target string) error {
	// 确保参数合理
	if domain == "" || target == "" {
		return fmt.Errorf("域名和目标地址不能为空")
	}

	// 确保目标地址格式正确
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		return c.core.PrintAndReturnError("目标地址格式不正确，必须以 http:// 或 https:// 开头")
	}

	// 确保配置目录和 Caddyfile 存在
	if err := c.EnsureConfigDir(); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	proxyConfig, err := c.RenderProxyConfig(domain, target)
	if err != nil {
		return fmt.Errorf("渲染代理配置失败: %v", err)
	}

	// 使用域名作为配置文件名获取项目配置路径
	configPath := c.GetProjectConfigPath(domain)

	// 写入代理配置
	if err := c.WriteConfig(configPath, proxyConfig); err != nil {
		return fmt.Errorf("写入代理配置失败: %v", err)
	}

	c.core.Infoln("添加反向代理配置成功")
	c.core.Infoln(fmt.Sprintf("代理配置文件: %s", configPath))

	// 重新加载 Caddy 使配置生效
	if err := c.Reload(); err != nil {
		return c.core.PrintAndReturnError("重新加载 Caddy 失败")
	}

	return nil
}

// RemoveProxy 移除指定域名的代理配置
func (c *Caddy) RemoveProxy(domain string) error {
	configPath := c.GetProjectConfigPath(domain)

	// 删除代理配置文件
	if err := c.RemoveConfig(configPath); err != nil {
		return fmt.Errorf("删除代理配置失败: %v", err)
	}

	// 重新加载 Caddy 使配置生效
	if err := c.Reload(); err != nil {
		return fmt.Errorf("移除代理后重新加载 Caddy 失败: %v", err)
	}

	return nil
}
