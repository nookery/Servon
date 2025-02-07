package software

// Plugin 定义软件插件接口
type Plugin interface {
	// Init 初始化插件
	Init() error
	// Name 返回插件名称
	Name() string
	// Register 注册插件提供的软件
	Register()
}

// RegisterPlugin 注册一个插件
func RegisterPlugin(p Plugin) error {
	if err := p.Init(); err != nil {
		return err
	}
	p.Register()
	return nil
}
