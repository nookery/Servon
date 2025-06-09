package core

import (
	"path/filepath"

	"servon/components"
	"servon/components/events"
	"servon/components/logger"
	"servon/core/contract"
	"servon/core/managers"
	"servon/core/providers"
	"servon/core/repo"
)

type App struct {
	eventBus events.IEventBus

	*providers.ManagerProvider
	*providers.CommandProvider

	AppLogger      *logger.LogUtil
	PluginRegistry *repo.PluginRegistry
	PluginManager  *managers.PluginManager
}

// New 创建App实例
func New() *App {
	eventBus := components.EventBus

	manager := managers.NewManager(eventBus)
	pluginRegistry := repo.NewPluginRegistry()

	app := &App{
		eventBus:        eventBus,
		ManagerProvider: providers.NewManagerProvider(eventBus, manager),
		CommandProvider: providers.NewCommandProvider(manager),
		AppLogger:       logger.NewLogUtil(filepath.Join(DataRootFolder, "logs")),
		PluginRegistry:  pluginRegistry,
		PluginManager:   managers.NewPluginManager(pluginRegistry),
	}

	return app
}

// RegisterWebPlugin 注册Web插件
func (app *App) RegisterWebPlugin(plugin contract.SuperWebPlugin) {
	app.PluginRegistry.RegisterWebPlugin(plugin)
}

// RegisterCLIPlugin 注册CLI插件
func (app *App) RegisterCLIPlugin(plugin contract.SuperCLIPlugin) {
	app.PluginRegistry.RegisterCLIPlugin(plugin)
}

// RegisterPlugin 注册完整插件
func (app *App) RegisterPlugin(plugin contract.SuperPlugin) {
	app.PluginRegistry.RegisterPlugin(plugin)
}

// GetWebPlugins 获取所有Web插件
func (app *App) GetWebPlugins() []contract.SuperWebPlugin {
	return app.PluginRegistry.GetWebPlugins()
}

// GetEnabledWebPlugins 获取所有启用的Web插件
func (app *App) GetEnabledWebPlugins() []contract.SuperWebPlugin {
	return app.PluginRegistry.GetEnabledWebPlugins()
}

// GetPluginByName 根据名称获取插件
func (app *App) GetPluginByName(name string) contract.SuperWebPlugin {
	return app.PluginRegistry.GetPluginByName(name)
}

// ListPlugins 列出所有插件信息
func (app *App) ListPlugins() []contract.PluginInfo {
	return app.PluginManager.ListPlugins()
}

// GetPluginInfo 获取指定插件的详细信息
func (app *App) GetPluginInfo(name string) (*contract.PluginInfo, error) {
	return app.PluginManager.GetPluginInfo(name)
}

// GetPluginStats 获取插件统计信息
func (app *App) GetPluginStats() map[string]interface{} {
	return app.PluginManager.GetPluginStats()
}
