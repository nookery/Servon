package web

import (
	"fmt"
	"net/http"

	"servon/internal/system"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	port   int
}

func NewServer(port int) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// 加载HTML模板
	router.LoadHTMLGlob("web/templates/*")

	// 设置静态文件目录
	router.Static("/static", "./web/static")

	return &Server{
		router: router,
		port:   port,
	}
}

func (s *Server) setupRoutes() {
	// API 路由
	api := s.router.Group("/api")
	{
		api.GET("/system/info", s.handleSystemInfo)
		// TODO: 添加更多API路由
	}

	// Web页面路由
	s.router.GET("/", s.handleHome)
	s.router.GET("/system", s.handleSystemPage)
}

func (s *Server) Start() error {
	s.setupRoutes()
	return s.router.Run(fmt.Sprintf(":%d", s.port))
}

// API处理函数
func (s *Server) handleSystemInfo(c *gin.Context) {
	info, err := system.GetSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

// 页面处理函数
func (s *Server) handleHome(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Servon - 服务器管理面板",
	})
}

func (s *Server) handleSystemPage(c *gin.Context) {
	info, _ := system.GetSystemInfo()
	c.HTML(http.StatusOK, "system.html", gin.H{
		"title": "系统信息",
		"info":  info,
	})
}
