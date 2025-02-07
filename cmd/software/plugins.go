package software

import "sync"

import "servon/cmd/contract"

var (
	plugins     = make(map[string]contract.Plugin)
	pluginMutex sync.RWMutex
)

// LoadPlugins 加载所有已注册的插件
func LoadPlugins() error {
	pluginMutex.RLock()
	defer pluginMutex.RUnlock()

	for _, plugin := range plugins {
		if err := contract.RegisterPlugin(plugin); err != nil {
			return err
		}
	}
	return nil
}

// RegisterBuiltinPlugins 注册内置插件
func RegisterBuiltinPlugins() {
	// 这里可以注册内置的插件
	// 例如：RegisterPlugin(&CaddyPlugin{})
}
