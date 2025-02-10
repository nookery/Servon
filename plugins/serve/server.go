package serve

import (
	"fmt"
	"servon/core"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	host   string
	port   int
	withUI bool
	*core.Core
}

func (p *ServePlugin) NewServer(host string, port int, withUI bool) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	router.HandleMethodNotAllowed = false
	router.SetTrustedProxies(nil)
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/favicon.ico"}}))
	router.Use(gin.Recovery())

	// 配置 CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 添加请求日志
	router.Use(func(c *gin.Context) {
		fmt.Printf("[%s] %s\n", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	return &Server{
		router: router,
		host:   host,
		port:   port,
		withUI: withUI,
		Core:   p.Core,
	}
}

func (s *Server) setupRoutes() {
	// 设置API路由
	s.setupAPIRoutes()

	// 如果启用了UI，设置UI路由
	if s.withUI {
		s.setupUIRoutes()
	}
}

func (s *Server) Start() error {
	s.setupRoutes()
	return s.router.Run(fmt.Sprintf("%s:%d", s.host, s.port))
}
