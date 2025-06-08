## 概述

Servon采用基于[Cobra](https://github.com/spf13/cobra)的插件系统架构，允许通过添加新的命令模块来扩展系统功能。每个插件都是一个独立的命令模块，可以为系统添加新的功能和命令。

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

### 1. 创建插件目录

在 `plugins` 目录下创建新的插件目录：

```bash
mkdir plugins/your_plugin
```

### 2. 实现插件入口

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

### 3. 实现具体命令

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

### 4. 注册插件

在 `plugins/root.go` 中注册你的插件：

```go
func init() {
    rootCmd.AddCommand(your_plugin.YourPluginCmd)
}
```

## 插件开发最佳实践

1. **模块化设计**
   - 将插件功能分解为独立的子命令
   - 保持每个命令的单一职责
   - 适当使用子命令组织复杂功能

2. **错误处理**
   - 提供清晰的错误信息
   - 使用适当的退出码
   - 实现优雅的错误恢复

3. **文档和帮助**
   - 为每个命令提供简短描述
   - 编写详细的使用说明
   - 包含使用示例

4. **配置管理**
   - 使用标准的配置管理方式
   - 支持命令行参数
   - 提供合理的默认值

5. **测试**
   - 编写单元测试
   - 包含集成测试
   - 测试边界情况

## 示例插件

以下是一个简单的示例插件实现：

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

1. **命令未注册**
   - 检查插件是否正确注册到根命令
   - 确认插件目录结构正确

2. **参数解析错误**
   - 检查参数验证逻辑
   - 确认命令使用说明正确

3. **依赖问题**
   - 确保所有依赖都已正确安装
   - 检查版本兼容性
