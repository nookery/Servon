package libs

import (
	"io"

	"github.com/gin-gonic/gin"
)

func NewWebServer(host string, port int, withUI bool) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	router.HandleMethodNotAllowed = false
	router.SetTrustedProxies(nil)
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/favicon.ico"},
		Output:    io.Discard,
	}))
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
	// router.Use(func(c *gin.Context) {
	// 	DefaultPrinter.Printf("[%s] %s\n", c.Request.Method, c.Request.URL.Path)
	// 	c.Next()
	// })

	return router
}
