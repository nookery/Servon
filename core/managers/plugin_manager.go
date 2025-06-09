package managers

import (
	"fmt"
	"servon/core/contract"
	"servon/core/repo"
)

// PluginManager 插件管理器，提供插件管理功能
type PluginManager struct {
	registry *repo.PluginRegistry
}

// NewPluginManager 创建新的插件管理器实例
func NewPluginManager(registry *repo.PluginRegistry) *PluginManager {
	return &PluginManager{
		registry: registry,
	}
}

// ListPlugins 列出所有插件信息
func (pm *PluginManager) ListPlugins() []contract.PluginInfo {
	var pluginInfos []contract.PluginInfo

	// 获取所有Web插件
	webPlugins := pm.registry.GetWebPlugins()
	for _, plugin := range webPlugins {
		info := contract.PluginInfo{
			Name:        plugin.GetName(),
			Description: plugin.GetDescription(),
			Enabled:     plugin.IsEnabled(),
			HasWebAPI:   true,
			APIPrefix:   plugin.GetAPIPrefix(),
		}
		pluginInfos = append(pluginInfos, info)
	}

	// 获取所有CLI插件（排除已经在Web插件中的）
	cliPlugins := pm.registry.GetCLIPlugins()
	for _, plugin := range cliPlugins {
		// 检查是否已经在Web插件中
		alreadyExists := false
		for _, existing := range pluginInfos {
			if existing.Name == plugin.GetName() {
				// 更新现有插件信息，标记为同时支持CLI
				existing.HasCLI = true
				alreadyExists = true
				break
			}
		}

		if !alreadyExists {
			info := contract.PluginInfo{
				Name:        plugin.GetName(),
				Description: plugin.GetDescription(),
				Enabled:     plugin.IsEnabled(),
				HasCLI:      true,
			}
			pluginInfos = append(pluginInfos, info)
		}
	}

	return pluginInfos
}

// GetPluginInfo 获取指定插件的详细信息
func (pm *PluginManager) GetPluginInfo(name string) (*contract.PluginInfo, error) {
	plugins := pm.ListPlugins()
	for _, plugin := range plugins {
		if plugin.Name == name {
			return &plugin, nil
		}
	}
	return nil, fmt.Errorf("plugin '%s' not found", name)
}

// GetEnabledPluginCount 获取启用的插件数量
func (pm *PluginManager) GetEnabledPluginCount() int {
	plugins := pm.ListPlugins()
	count := 0
	for _, plugin := range plugins {
		if plugin.Enabled {
			count++
		}
	}
	return count
}

// GetPluginStats 获取插件统计信息
func (pm *PluginManager) GetPluginStats() map[string]interface{} {
	plugins := pm.ListPlugins()
	webCount := 0
	cliCount := 0
	enabledCount := 0

	for _, plugin := range plugins {
		if plugin.HasWebAPI {
			webCount++
		}
		if plugin.HasCLI {
			cliCount++
		}
		if plugin.Enabled {
			enabledCount++
		}
	}

	return map[string]interface{}{
		"total_plugins":   len(plugins),
		"web_plugins":     webCount,
		"cli_plugins":     cliCount,
		"enabled_plugins": enabledCount,
	}
}
