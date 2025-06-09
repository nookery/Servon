package web_server

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// SetPort 设置端口号
func (ws *WebServer) SetPort(port int) {
	ws.config.Port = port
}

// SetHost 设置主机
func (ws *WebServer) SetHost(host string) {
	ws.config.Host = host
}

// SetVerbose 设置详细日志模式
func (ws *WebServer) SetVerbose(verbose bool) {
	ws.config.Verbose = verbose
}

// SetLogger 设置自定义日志器
func (ws *WebServer) SetLogger(logger Logger) {
	if logger != nil {
		ws.logger = logger
		ws.config.Logger = logger
	} else {
		ws.logger = NewDefaultLogger()
		ws.config.Logger = nil
	}
}

// SetConfig 设置服务器配置
func (ws *WebServer) SetConfig(config WebServerConfig) {
	ws.config = config
}

// SetupCORS 设置CORS
func (ws *WebServer) SetupCORS() {
	ws.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// SetupLogging 设置日志中间件
func (ws *WebServer) SetupLogging() {
	ws.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
}

// HealthCheck 健康检查端点
func (ws *WebServer) SetupHealthCheck() {
	ws.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})
}

// SetupStaticFiles 设置静态文件服务
func (ws *WebServer) SetupStaticFiles(relativePath, root string) {
	ws.Static(relativePath, root)
}

// SetupTemplates 设置模板
func (ws *WebServer) SetupTemplates(pattern string) {
	ws.LoadHTMLGlob(pattern)
}
