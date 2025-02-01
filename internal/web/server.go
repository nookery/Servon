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

	return &Server{
		router: router,
		port:   port,
	}
}

func (s *Server) setupRoutes() {
	// API 路由
	api := s.router.Group("/api")
	{
		api.GET("/system/basic", s.handleBasicInfo) // 新增基本信息接口
		// TODO: 添加更多API路由
	}
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

// 新增：处理基本系统信息的接口
func (s *Server) handleBasicInfo(c *gin.Context) {
	info, err := system.GetBasicSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}
