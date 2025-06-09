package repo

import (
	"servon/core/contract"
	"sync"
)

// PluginRegistry 插件注册表，用于管理所有插件
type PluginRegistry struct {
	mu         sync.RWMutex
	webPlugins []contract.SuperWebPlugin
	cliPlugins []contract.SuperCLIPlugin
	plugins    []contract.SuperPlugin
}

// NewPluginRegistry 创建新的插件注册表实例
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		webPlugins: make([]contract.SuperWebPlugin, 0),
		cliPlugins: make([]contract.SuperCLIPlugin, 0),
		plugins:    make([]contract.SuperPlugin, 0),
	}
}

// RegisterWebPlugin 注册Web插件
func (r *PluginRegistry) RegisterWebPlugin(plugin contract.SuperWebPlugin) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.webPlugins = append(r.webPlugins, plugin)
}

// RegisterCLIPlugin 注册CLI插件
func (r *PluginRegistry) RegisterCLIPlugin(plugin contract.SuperCLIPlugin) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cliPlugins = append(r.cliPlugins, plugin)
}

// RegisterPlugin 注册完整插件（同时支持Web和CLI）
func (r *PluginRegistry) RegisterPlugin(plugin contract.SuperPlugin) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plugins = append(r.plugins, plugin)
}

// GetWebPlugins 获取所有Web插件
func (r *PluginRegistry) GetWebPlugins() []contract.SuperWebPlugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 合并独立的Web插件和完整插件中的Web功能
	allWebPlugins := make([]contract.SuperWebPlugin, 0, len(r.webPlugins)+len(r.plugins))
	allWebPlugins = append(allWebPlugins, r.webPlugins...)
	for _, plugin := range r.plugins {
		allWebPlugins = append(allWebPlugins, plugin)
	}
	return allWebPlugins
}

// GetCLIPlugins 获取所有CLI插件
func (r *PluginRegistry) GetCLIPlugins() []contract.SuperCLIPlugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 合并独立的CLI插件和完整插件中的CLI功能
	allCLIPlugins := make([]contract.SuperCLIPlugin, 0, len(r.cliPlugins)+len(r.plugins))
	allCLIPlugins = append(allCLIPlugins, r.cliPlugins...)
	for _, plugin := range r.plugins {
		allCLIPlugins = append(allCLIPlugins, plugin)
	}
	return allCLIPlugins
}

// GetAllPlugins 获取所有完整插件
func (r *PluginRegistry) GetAllPlugins() []contract.SuperPlugin {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]contract.SuperPlugin(nil), r.plugins...)
}

// GetPluginByName 根据名称获取插件
func (r *PluginRegistry) GetPluginByName(name string) contract.SuperWebPlugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 在Web插件中查找
	for _, plugin := range r.webPlugins {
		if plugin.GetName() == name {
			return plugin
		}
	}

	// 在完整插件中查找
	for _, plugin := range r.plugins {
		if plugin.GetName() == name {
			return plugin
		}
	}

	return nil
}

// GetPluginCount 获取插件总数
func (r *PluginRegistry) GetPluginCount() (webCount, cliCount, totalCount int) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	webCount = len(r.webPlugins) + len(r.plugins)
	cliCount = len(r.cliPlugins) + len(r.plugins)
	totalCount = len(r.webPlugins) + len(r.cliPlugins) + len(r.plugins)

	return
}

// GetEnabledWebPlugins 获取所有启用的Web插件
func (r *PluginRegistry) GetEnabledWebPlugins() []contract.SuperWebPlugin {
	allPlugins := r.GetWebPlugins()
	enabledPlugins := make([]contract.SuperWebPlugin, 0)

	for _, plugin := range allPlugins {
		if plugin.IsEnabled() {
			enabledPlugins = append(enabledPlugins, plugin)
		}
	}

	return enabledPlugins
}
