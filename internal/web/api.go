package web

import (
	"servon/internal/web/handler"
)

// setupAPIRoutes 设置所有API路由
func (s *Server) setupAPIRoutes() {
	h := handler.New()
	sh := handler.NewSoftwareHandler()
	api := s.router.Group("/web_api")
	{
		api.GET("/system/resources", h.HandleSystemResources)
		api.GET("/system/network", h.HandleNetworkResources)
		api.GET("/system/user", h.HandleCurrentUser)
		api.GET("/system/os", h.HandleOSInfo)
		api.GET("/system/basic", h.HandleBasicInfo)
		api.GET("/system/software", sh.HandleGetSoftwareList)
		api.GET("/system/software/:name/install", sh.HandleInstallSoftware)
		api.GET("/system/software/:name/uninstall", sh.HandleUninstallSoftware)
		api.POST("/system/software/:name/stop", sh.HandleStopSoftware)
		api.GET("/system/software/:name/status", sh.HandleGetSoftwareStatus)
		api.GET("/system/processes", h.HandleProcessList)
		api.GET("/system/files", h.HandleFileList)
		api.GET("/system/ports", h.HandlePortList)
	}
}
