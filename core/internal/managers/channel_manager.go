package managers

var DefaultChannelManager = newChannelManager()

type ChannelManager struct {
	channels map[string]chan string
}

func newChannelManager() *ChannelManager {
	return &ChannelManager{
		channels: make(map[string]chan string),
	}
}

func (c *ChannelManager) GetChannel(name string) chan string {
	return c.channels[name]
}
