package libs

type ChannelManager struct {
	channels map[string]chan string
}

func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		channels: make(map[string]chan string),
	}
}

func (c *ChannelManager) GetChannel(name string) chan string {
	return c.channels[name]
}
