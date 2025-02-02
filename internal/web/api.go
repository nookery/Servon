package web

import (
	"servon/internal/web/handler"
)

// setupAPIRoutes 设置所有API路由
func (s *Server) setupAPIRoutes() {
	h := handler.New()
	api := s.router.Group("/web_api")
	{
		api.GET("/system/resources", h.HandleSystemResources)
		api.GET("/system/user", h.HandleCurrentUser)
		api.GET("/system/os", h.HandleOSInfo)
		api.GET("/system/basic", h.HandleBasicInfo)
		api.GET("/system/software", h.HandleSoftwareList)
		api.GET("/system/software/:name/install", h.HandleSoftwareInstall)
		api.GET("/system/software/:name/uninstall", h.HandleSoftwareUninstall)
		api.POST("/system/software/:name/stop", h.HandleSoftwareStop)
		api.GET("/system/software/:name/status", h.HandleSoftwareStatus)
		api.GET("/system/processes", h.HandleProcessList)
		api.GET("/system/files", h.HandleFileList)
		api.GET("/system/ports", h.HandlePortList)
	}
}
