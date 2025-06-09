# Web Server Logger 使用说明

## 概述

`web_server` 组件现在支持内部默认 logger 和外部自定义 logger。默认情况下使用内部 logger，同时允许用户提供自定义 logger 实现。

## Logger 接口

所有 logger 必须实现以下接口：

```go
type Logger interface {
    Infof(format string, args ...interface{})
    Errorf(format string, args ...interface{})
    Warnf(format string, args ...interface{})
    Debugf(format string, args ...interface{})
}
```

## 内置 Logger 实现

### DefaultLogger

- 使用标准库 `log` 包
- 输出格式：`[INFO] YYYY/MM/DD HH:MM:SS 消息内容`
- 支持不同日志级别的颜色输出

### NoOpLogger

- 空操作 logger，不输出任何内容
- 用于禁用日志输出的场景

## 使用方式

### 1. 使用默认内部 Logger

```go
config := web_server.WebServerConfig{
    Host:    "127.0.0.1",
    Port:    8080,
    Verbose: true,
    // Logger 字段为 nil，将自动使用 DefaultLogger
}
ws := web_server.NewWebServer(config)
```

### 2. 使用自定义 Logger

```go
// 实现自定义 logger
type CustomLogger struct {
    logger *log.Logger
}

func (c *CustomLogger) Infof(format string, args ...interface{}) {
    c.logger.Printf("[INFO] "+format, args...)
}

// ... 实现其他方法

// 使用自定义 logger
customLogger := &CustomLogger{
    logger: log.New(os.Stdout, "[CUSTOM] ", log.LstdFlags),
}

config := web_server.WebServerConfig{
    Host:    "127.0.0.1",
    Port:    8080,
    Verbose: true,
    Logger:  customLogger, // 使用自定义 logger
}
ws := web_server.NewWebServer(config)
```

### 3. 运行时更换 Logger

```go
// 获取当前 logger
currentLogger := ws.GetLogger()

// 设置新的 logger
ws.SetLogger(newLogger)

// 设置为 nil 将使用默认 logger
ws.SetLogger(nil)
```

## 日志输出示例

### 默认 Logger 输出

```bash
[INFO] 2025/06/09 09:10:05 🚀 启动服务器在端口 8080
[INFO] 2025/06/09 09:10:05 ✅ 服务器已成功关闭
```

### 自定义 Logger 输出

```bash
[CUSTOM] 2025/06/09 09:11:07 [INFO] 这是来自自定义 logger 的信息
[CUSTOM] 2025/06/09 09:11:07 [ERROR] 这是来自自定义 logger 的错误
```

## 配置选项

- `Verbose`: 设置为 `true` 启用详细日志输出
- `Logger`: 自定义 logger 实例，为 `nil` 时使用默认 logger

## 最佳实践

1. **生产环境**：建议使用自定义 logger，集成到现有的日志系统中
2. **开发环境**：可以使用默认 logger，便于快速调试
3. **测试环境**：可以使用 `NoOpLogger` 禁用日志输出
4. **日志级别**：根据 `Verbose` 配置控制日志详细程度

## 示例代码

完整的使用示例请参考 `examples/custom_logger_example.go` 文件。
