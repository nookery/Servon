package providers

import (
	"servon/components/web_server"
	"servon/core/managers"
	"servon/core/web/routers"
)

type WebProvider struct {
	Server *web_server.WebServer
}

func NewWebProvider(manager *managers.FullManager) *WebProvider {
	server := web_server.NewWebServer()
	webProvider := &WebProvider{
		Server: server,
	}

	routers.Setup(manager, server.Engine, true)

	return webProvider
}
