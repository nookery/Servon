package routers

import (
	"servon/core"
	"servon/plugins/web/handlers"
)

func Setup(app *core.App, isDev bool) {
	app.PrintInfof("Setup Web Router IsDev: %t", isDev)
	handler := &handlers.WebHandler{
		App: app,
	}

	router := &WebRouter{
		Handler: handler,
		App:     app,
	}

	router.setup(isDev)
}

type WebRouter struct {
	Handler *handlers.WebHandler
	*core.App
}

func (w *WebRouter) setup(isDev bool) {
	api := w.GetRouter().Group("/web_api")

	w.SetupSoftRouter(api)
	w.SetupProcessRouter(api)
	w.SetupInfoRouter(api)
	w.SetupGitHubRouter(api)
	w.SetupTaskRouter(api)
	w.SetupCronRouter(api)
	w.SetupFileRouter(api)
	w.SetupPortRouter(api)
	w.SetupLogsRouter(api)
	w.SetupUserRouter(api)
	w.SetupWebApiRouter(api)
	w.SetupDeployRouter(api)

	w.PrintSuccess("API Router Setup Success")
	w.SetupUIRoutes(w.GetRouter(), isDev)
}
