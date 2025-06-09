package web_server

import (
	"fmt"
	"log"
	"os"
)

// CustomLogger 实现自定义 logger
type CustomLogger struct {
	logger *log.Logger
}

// NewCustomLogger 创建自定义 logger
func NewCustomLogger() *CustomLogger {
	return &CustomLogger{
		logger: log.New(os.Stdout, "[CUSTOM] ", log.LstdFlags),
	}
}

// Infof 实现 Logger 接口的 Infof 方法
func (c *CustomLogger) Infof(format string, args ...interface{}) {
	c.logger.Printf("[INFO] "+format, args...)
}

// Errorf 实现 Logger 接口的 Errorf 方法
func (c *CustomLogger) Errorf(format string, args ...interface{}) {
	c.logger.Printf("[ERROR] "+format, args...)
}

// Warnf 实现 Logger 接口的 Warnf 方法
func (c *CustomLogger) Warnf(format string, args ...interface{}) {
	c.logger.Printf("[WARN] "+format, args...)
}

// Debugf 实现 Logger 接口的 Debugf 方法
func (c *CustomLogger) Debugf(format string, args ...interface{}) {
	c.logger.Printf("[DEBUG] "+format, args...)
}

func Start() {
	fmt.Println("=== 使用默认内部 Logger ===")
	// 使用默认内部 logger
	config1 := WebServerConfig{
		Host:    "127.0.0.1",
		Port:    8081,
		Verbose: true,
		// Logger 字段为 nil，将使用默认内部 logger
	}
	ws1 := NewWebServer()
	ws1.SetConfig(config1)
	fmt.Printf("WebServer 1 使用的 Logger 类型: %T\n", ws1.GetLogger())

	fmt.Println("\n=== 使用自定义 Logger ===")
	// 使用自定义 logger
	customLogger := NewCustomLogger()
	config2 := WebServerConfig{
		Host:    "127.0.0.1",
		Port:    8082,
		Verbose: true,
		Logger:  customLogger, // 使用自定义 logger
	}
	ws2 := NewWebServer()
	ws2.SetConfig(config2)
	fmt.Printf("WebServer 2 使用的 Logger 类型: %T\n", ws2.GetLogger())

	fmt.Println("\n=== 运行时更换 Logger ===")
	// 运行时更换 logger
	ws1.SetLogger(customLogger)
	fmt.Printf("WebServer 1 更换后的 Logger 类型: %T\n", ws1.GetLogger())

	fmt.Println("\n=== 测试 Logger 输出 ===")
	// 测试 logger 输出
	ws1.GetLogger().Infof("这是来自 WebServer 1 的信息日志")
	ws2.GetLogger().Infof("这是来自 WebServer 2 的信息日志")
	ws1.GetLogger().Errorf("这是来自 WebServer 1 的错误日志")
	ws2.GetLogger().Warnf("这是来自 WebServer 2 的警告日志")
}
