## 概述

Servon采用基于[Cobra](https://github.com/spf13/cobra)的插件系统架构，允许通过添加新的命令模块来扩展系统功能。每个插件都是一个独立的命令模块，可以为系统添加新的功能和命令。

从v2.0开始，Servon引入了全新的插件注册机制，提供了统一的插件管理和发现功能，支持Web API和CLI功能的动态注册。

## Servon v2.0+ 新插件注册机制

### 概述

Servon Core 提供了一套完整的插件注册机制，允许插件动态注册Web API和CLI功能，实现了插件的统一管理和发现。

### 核心特性

- **统一接口**: 通过标准化的接口定义，实现插件的统一管理
- **动态发现**: Server插件可以自动发现并注册所有Web插件的路由
- **类型安全**: 通过接口约束确保插件实现的一致性
- **解耦设计**: 插件与核心系统解耦，便于独立开发和测试
- **扩展性**: 易于添加新的插件类型和功能

### 核心组件

#### 1. 插件接口 (contract/plugin.go)

定义了三个核心接口：

- **WebPlugin**: 定义Web功能插件接口
- **CLIPlugin**: 定义CLI功能插件接口  
- **Plugin**: 组合接口，同时支持Web和CLI功能

```go
type WebPlugin interface {
    GetName() string
    SetupWebRoutes(router *gin.RouterGroup, manager interface{})
    GetAPIPrefix() string
    GetDescription() string
    IsEnabled() bool
}

type CLIPlugin interface {
    GetName() string
    GetDescription() string
    IsEnabled() bool
}

type Plugin interface {
    WebPlugin
    CLIPlugin
}
```

#### 2. 插件注册表 (plugin_registry.go)

负责管理所有插件的注册和获取：

- `RegisterWebPlugin()`: 注册Web插件
- `RegisterCLIPlugin()`: 注册CLI插件
- `RegisterPlugin()`: 注册完整插件
- `GetWebPlugins()`: 获取所有Web插件
- `GetEnabledWebPlugins()`: 获取启用的Web插件
- `GetPluginByName()`: 根据名称查找插件

#### 3. 插件管理器 (plugin_manager.go)

提供插件管理功能：

- `ListPlugins()`: 列出所有插件信息
- `GetPluginInfo()`: 获取指定插件详情
- `GetPluginStats()`: 获取插件统计信息

#### 4. 应用程序集成 (app.go)

App结构体集成了插件注册机制，提供统一的插件管理接口。

### 插件注册方式

在插件的Setup函数中，可以使用以下方法注册插件：

```go
func Setup(app *core.App) {
    plugin := &MyPlugin{
        name:        "myplugin",
        description: "My awesome plugin",
        enabled:     true,
    }
    
    // 注册Web插件
    app.RegisterWebPlugin(plugin)
    
    // 注册CLI插件
    app.RegisterCLIPlugin(plugin)
    
    // 或者注册完整插件（同时支持Web和CLI）
    app.RegisterPlugin(plugin)
}
```

### 插件管理API

新的插件注册机制提供了丰富的插件管理API：

```go
// 获取所有插件信息
plugins := app.ListPlugins()

// 获取特定插件信息
info, err := app.GetPluginInfo("pluginName")

// 获取插件统计信息
stats := app.GetPluginStats()

// 获取启用的Web插件
enabledWebPlugins := app.GetEnabledWebPlugins()

// 根据名称查找插件
plugin := app.GetPluginByName("pluginName")
```

### 架构优势

1. **解耦**: 插件与核心系统解耦，便于独立开发和测试
2. **动态发现**: Server插件可以动态发现并注册所有Web插件的路由
3. **统一管理**: 提供统一的插件管理接口
4. **类型安全**: 通过接口定义确保插件实现的一致性
5. **扩展性**: 易于添加新的插件类型和功能

### 注意事项

1. 插件必须实现相应的接口才能被正确注册
2. Web插件的SetupWebRoutes方法中的manager参数使用interface{}类型以避免循环导入
3. 插件名称应该唯一，避免冲突
4. 插件的IsEnabled()方法决定插件是否被激活使用

## 插件结构

### 目录结构

插件文件统一存放在 `plugins` 目录下，每个插件都有自己的子目录：

```
plugins/
├── astro/          # Astro构建工具插件
├── bt/             # 宝塔面板集成插件
├── caddy/          # Caddy服务器插件
├── clash/          # Clash代理插件
├── database/       # 数据库管理插件
├── git/            # Git操作插件
├── github_runner/  # GitHub Actions Runner插件
├── nodejs/         # Node.js环境插件
├── npm/            # NPM包管理插件
├── pm2/            # PM2进程管理插件
├── supervisor/     # Supervisor进程管理插件
├── system/         # 系统信息插件
└── xcode/          # Xcode开发工具插件
```

### 插件基本结构

一个典型的插件包含以下文件：

```go
plugins/example/
├── root.go         # 插件入口点和主要命令定义
├── commands.go     # 具体命令实现
└── types.go        # 插件使用的数据类型定义
```

## 开发新插件

### 方式一：使用新的插件注册机制 (推荐)

#### 1. 创建插件目录

在 `plugins` 目录下创建新的插件目录：

```bash
mkdir plugins/your_plugin
```

#### 2. 实现插件结构体

创建插件主文件 `your_plugin.go`：

```go
package your_plugin

import (
    "servon/core"
    "servon/core/contract"
    "github.com/gin-gonic/gin"
    "github.com/spf13/cobra"
)

// YourPlugin 插件结构体
type YourPlugin struct {
    name        string
    description string
    enabled     bool
}

// GetName 返回插件名称
func (p *YourPlugin) GetName() string {
    return p.name
}

// GetDescription 返回插件描述
func (p *YourPlugin) GetDescription() string {
    return p.description
}

// IsEnabled 检查插件是否启用
func (p *YourPlugin) IsEnabled() bool {
    return p.enabled
}

// GetAPIPrefix 返回API前缀路径
func (p *YourPlugin) GetAPIPrefix() string {
    return "/your-plugin"
}

// SetupWebRoutes 设置Web路由
func (p *YourPlugin) SetupWebRoutes(router *gin.RouterGroup, manager interface{}) {
    // 设置API路由
    router.GET("/info", p.handleInfo)
    router.POST("/action", p.handleAction)
}

// handleInfo 处理信息请求
func (p *YourPlugin) handleInfo(c *gin.Context) {
    c.JSON(200, gin.H{
        "name":        p.GetName(),
        "description": p.GetDescription(),
        "status":      "running",
    })
}

// handleAction 处理操作请求
func (p *YourPlugin) handleAction(c *gin.Context) {
    // 实现具体的操作逻辑
    c.JSON(200, gin.H{"message": "Action completed"})
}
```

#### 3. 实现CLI命令 (可选)

创建 `root.go` 文件：

```go
package your_plugin

import (
    "github.com/spf13/cobra"
)

var YourPluginCmd = &cobra.Command{
    Use:   "your-plugin",
    Short: "Your plugin description",
    Long:  "A longer description of your plugin",
}

var infoCmd = &cobra.Command{
    Use:   "info",
    Short: "Show plugin information",
    Run: func(cmd *cobra.Command, args []string) {
        // CLI命令实现逻辑
    },
}

func init() {
    YourPluginCmd.AddCommand(infoCmd)
}
```

#### 4. 注册插件

创建 `setup.go` 文件：

```go
package your_plugin

import (
    "servon/core"
)

// Setup 注册插件到应用程序
func Setup(app *core.App) {
    // 创建插件实例
    plugin := &YourPlugin{
        name:        "your-plugin",
        description: "Your awesome plugin",
        enabled:     true,
    }
    
    // 注册为Web插件
    app.RegisterWebPlugin(plugin)
    
    // 如果同时支持CLI，注册CLI命令
    app.RootCmd.AddCommand(YourPluginCmd)
    
    // 或者注册为完整插件（需要实现完整的Plugin接口）
    // app.RegisterPlugin(plugin)
}
```

#### 5. 在main.go中导入插件

在 `main.go` 中添加插件导入：

```go
import (
    // ... 其他导入
    "servon/plugins/your_plugin"
)

func main() {
    app := core.New()
    
    // ... 其他插件设置
    your_plugin.Setup(app)
    
    app.Execute()
}
```

### 方式二：传统Cobra命令方式

#### 1. 创建插件目录

```bash
mkdir plugins/your_plugin
```

#### 2. 实现插件入口

在插件目录中创建 `root.go` 文件：

```go
package your_plugin

import (
    "github.com/spf13/cobra"
)

var YourPluginCmd = &cobra.Command{
    Use:   "your-plugin",
    Short: "Your plugin description",
    Long:  "A longer description of your plugin",
}

func init() {
    // 添加子命令
    YourPluginCmd.AddCommand(yourSubCommand)
}
```

#### 3. 实现具体命令

创建命令实现文件：

```go
var yourSubCommand = &cobra.Command{
    Use:   "subcommand",
    Short: "Subcommand description",
    Run: func(cmd *cobra.Command, args []string) {
        // 命令实现逻辑
    },
}
```

#### 4. 注册插件

在 `main.go` 中注册你的插件：

```go
import "servon/plugins/your_plugin"

func main() {
    app := core.New()
    app.RootCmd.AddCommand(your_plugin.YourPluginCmd)
    app.Execute()
}
```

## 插件开发最佳实践

### 新插件注册机制最佳实践

1. **接口实现**
   - 完整实现所需的插件接口方法
   - 提供有意义的插件名称和描述
   - 正确实现IsEnabled()方法控制插件状态

2. **Web路由设计**
   - 使用有意义的API前缀路径
   - 遵循RESTful API设计原则
   - 提供适当的HTTP状态码和响应格式
   - 实现错误处理和参数验证

3. **插件结构**
   - 将插件逻辑封装在结构体中
   - 使用依赖注入获取所需的管理器
   - 保持插件代码的模块化和可测试性

4. **注册策略**
   - 根据插件功能选择合适的注册方式（Web、CLI或完整插件）
   - 在Setup函数中进行统一的插件注册
   - 避免在init函数中进行复杂的初始化操作

### 通用最佳实践

1. **模块化设计**
   - 将插件功能分解为独立的子命令
   - 保持每个命令的单一职责
   - 适当使用子命令组织复杂功能

2. **错误处理**
   - 提供清晰的错误信息
   - 使用适当的退出码
   - 实现优雅的错误恢复
   - 在Web API中返回标准的错误响应格式

3. **文档和帮助**
   - 为每个命令提供简短描述
   - 编写详细的使用说明
   - 包含使用示例
   - 为Web API提供OpenAPI文档

4. **配置管理**
   - 使用标准的配置管理方式
   - 支持命令行参数
   - 提供合理的默认值
   - 支持环境变量配置

5. **测试**
   - 编写单元测试
   - 包含集成测试
   - 测试边界情况
   - 为Web API编写API测试

6. **安全性**
   - 验证输入参数
   - 实现适当的权限控制
   - 避免敏感信息泄露
   - 使用HTTPS进行API通信

## 示例插件

### 新插件注册机制示例

以下是一个使用新插件注册机制的完整示例：

```go
// example_plugin.go
package example

import (
    "fmt"
    "servon/core"
    "servon/core/contract"
    "github.com/gin-gonic/gin"
    "github.com/spf13/cobra"
)

// ExamplePlugin 示例插件结构体
type ExamplePlugin struct {
    name        string
    description string
    enabled     bool
    version     string
}

// GetName 返回插件名称
func (p *ExamplePlugin) GetName() string {
    return p.name
}

// GetDescription 返回插件描述
func (p *ExamplePlugin) GetDescription() string {
    return p.description
}

// IsEnabled 检查插件是否启用
func (p *ExamplePlugin) IsEnabled() bool {
    return p.enabled
}

// GetAPIPrefix 返回API前缀路径
func (p *ExamplePlugin) GetAPIPrefix() string {
    return "/example"
}

// SetupWebRoutes 设置Web路由
func (p *ExamplePlugin) SetupWebRoutes(router *gin.RouterGroup, manager interface{}) {
    // 基础信息API
    router.GET("/info", p.handleInfo)
    router.GET("/version", p.handleVersion)
    
    // 功能API
    router.POST("/hello", p.handleHello)
    router.GET("/status", p.handleStatus)
}

// handleInfo 处理信息请求
func (p *ExamplePlugin) handleInfo(c *gin.Context) {
    c.JSON(200, gin.H{
        "name":        p.GetName(),
        "description": p.GetDescription(),
        "version":     p.version,
        "enabled":     p.IsEnabled(),
        "api_prefix":  p.GetAPIPrefix(),
    })
}

// handleVersion 处理版本请求
func (p *ExamplePlugin) handleVersion(c *gin.Context) {
    c.JSON(200, gin.H{
        "version": p.version,
        "plugin":  p.GetName(),
    })
}

// handleHello 处理问候请求
func (p *ExamplePlugin) handleHello(c *gin.Context) {
    var req struct {
        Name string `json:"name" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Name is required"})
        return
    }
    
    c.JSON(200, gin.H{
        "message": fmt.Sprintf("Hello, %s! Greetings from %s plugin.", req.Name, p.GetName()),
        "plugin":  p.GetName(),
    })
}

// handleStatus 处理状态请求
func (p *ExamplePlugin) handleStatus(c *gin.Context) {
    c.JSON(200, gin.H{
        "status":  "running",
        "plugin":  p.GetName(),
        "enabled": p.IsEnabled(),
    })
}
```

```go
// root.go - CLI命令实现
package example

import (
    "fmt"
    "github.com/spf13/cobra"
)

var ExampleCmd = &cobra.Command{
    Use:   "example",
    Short: "Example plugin demonstration",
    Long:  "A longer description of the example plugin with new registration mechanism",
}

var sayHelloCmd = &cobra.Command{
    Use:   "hello [name]",
    Short: "Say hello to someone",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        name := args[0]
        fmt.Printf("Hello, %s! Greetings from Example plugin.\n", name)
    },
}

var infoCmd = &cobra.Command{
    Use:   "info",
    Short: "Show plugin information",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Example Plugin Information:")
        fmt.Println("Name: example")
        fmt.Println("Description: Example plugin demonstration")
        fmt.Println("Version: 1.0.0")
        fmt.Println("API Prefix: /example")
    },
}

func init() {
    ExampleCmd.AddCommand(sayHelloCmd)
    ExampleCmd.AddCommand(infoCmd)
}
```

```go
// setup.go - 插件注册
package example

import (
    "servon/core"
)

// Setup 注册示例插件到应用程序
func Setup(app *core.App) {
    // 创建插件实例
    plugin := &ExamplePlugin{
        name:        "example",
        description: "Example plugin demonstration with new registration mechanism",
        enabled:     true,
        version:     "1.0.0",
    }
    
    // 注册为Web插件
    app.RegisterWebPlugin(plugin)
    
    // 注册CLI命令
    app.RootCmd.AddCommand(ExampleCmd)
    
    // 如果插件同时实现了CLIPlugin接口，也可以注册为完整插件
    // app.RegisterPlugin(plugin)
}
```

### 传统Cobra命令示例

以下是一个简单的传统插件实现：

```go
package example

import (
    "fmt"
    "github.com/spf13/cobra"
)

var ExampleCmd = &cobra.Command{
    Use:   "example",
    Short: "Example plugin demonstration",
    Long:  "A longer description of the example plugin",
}

var sayHelloCmd = &cobra.Command{
    Use:   "hello [name]",
    Short: "Say hello to someone",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        name := args[0]
        fmt.Printf("Hello, %s!\n", name)
    },
}

func init() {
    ExampleCmd.AddCommand(sayHelloCmd)
}
```

## 插件API参考

插件可以使用以下核心功能：

1. **命令注册**
   - `AddCommand()` - 添加子命令
   - `PersistentFlags()` - 添加持久标志
   - `Flags()` - 添加本地标志

2. **参数处理**
   - `Args` - 参数验证
   - `PreRun` - 预运行钩子
   - `Run` - 主要命令逻辑
   - `PostRun` - 后运行钩子

3. **输出控制**
   - 使用 `color` 包进行输出着色
   - 标准输出和错误输出控制
   - 进度显示和状态更新

## 调试插件

1. **启用调试日志**
   ```bash
   export DEBUG=1
   ```

2. **使用 `-v` 标志获取详细输出**
   ```bash
   servon your-plugin -v
   ```

3. **检查命令结构**
   ```bash
   servon help your-plugin
   ```

## 常见问题

### 新插件注册机制相关问题

1. **插件未被发现**
   - 检查插件是否正确实现了相应的接口（WebPlugin、CLIPlugin或Plugin）
   - 确认在Setup函数中正确调用了注册方法
   - 验证插件的IsEnabled()方法返回true

2. **Web路由不工作**
   - 检查SetupWebRoutes方法是否正确实现
   - 确认API前缀路径设置正确
   - 验证Server插件是否正确集成了插件注册机制

3. **循环导入错误**
   - 避免在contract包中导入managers包
   - 使用interface{}类型作为manager参数
   - 在插件中进行类型断言获取具体的管理器

4. **插件信息不正确**
   - 检查GetName()、GetDescription()等方法的实现
   - 确认插件结构体字段正确初始化
   - 验证插件注册时传入的参数

### 通用问题

1. **命令未注册**
   - 检查插件是否正确注册到根命令
   - 确认插件目录结构正确
   - 验证main.go中是否正确导入和调用Setup函数

2. **参数解析错误**
   - 检查参数验证逻辑
   - 确认命令使用说明正确
   - 验证Cobra命令定义的Args字段

3. **依赖问题**
   - 确保所有依赖都已正确安装
   - 检查版本兼容性
   - 运行`go mod tidy`更新依赖

4. **编译错误**
   - 检查导入路径是否正确
   - 确认接口实现是否完整
   - 验证Go语法和类型匹配
