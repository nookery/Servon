# Server 插件

## 概述

Server 插件是 Servon 的核心插件之一，提供了一个功能强大的 Web 管理面板，让用户可以通过浏览器界面管理服务器的各种资源和服务。该插件基于 Gin 框架构建，提供了丰富的 RESTful API 和现代化的 Web 界面。

## 主要功能

### 🚀 服务器管理

- **启动服务器**: 支持后台运行，可配置端口、主机地址
- **停止服务器**: 优雅关闭服务器进程
- **重启服务器**: 无缝重启服务，保持服务连续性
- **开发模式**: 支持开发环境的热重启

### 📊 系统监控

- **系统资源监控**: CPU、内存、磁盘使用情况
- **网络资源监控**: 网络接口状态和流量统计
- **进程管理**: 查看和管理系统进程
- **端口监控**: 监控端口占用情况

### 🗂️ 文件管理

- **文件浏览**: 浏览服务器文件系统
- **文件操作**: 创建、删除、重命名、复制文件
- **文件编辑**: 在线编辑文件内容
- **文件下载**: 下载服务器文件到本地
- **批量操作**: 支持批量删除文件

### ⚙️ 服务管理

- **服务列表**: 查看所有后台服务
- **服务控制**: 启动、停止、重启服务
- **服务配置**: 管理服务配置文件
- **服务日志**: 查看服务运行日志
- **后台服务**: 添加和删除后台服务

### 📝 日志管理

- **日志文件列表**: 查看所有日志文件
- **日志内容读取**: 实时查看日志内容
- **日志搜索**: 搜索特定日志内容
- **日志统计**: 获取日志统计信息
- **日志清理**: 清理旧日志文件

### ⏰ 定时任务

- **任务管理**: 创建、更新、删除定时任务
- **任务调度**: 启用/禁用定时任务
- **任务监控**: 查看任务执行状态

### 🔧 软件管理

- **软件列表**: 查看已安装软件
- **软件安装**: 安装新软件包
- **软件卸载**: 卸载不需要的软件
- **软件控制**: 启动/停止软件服务

### 🌐 集成功能

- **GitHub 集成**: GitHub 仓库管理和 Webhook 处理
- **部署管理**: 自动化部署功能
- **拓扑管理**: 网关和项目拓扑管理
- **用户管理**: 用户账户管理

## 使用方法

### 基本命令

```bash
# 启动服务器
servon server start

# 指定端口和主机
servon server start --port 8080 --host 0.0.0.0

# 启用详细日志
servon server start --verbose

# 开发模式
servon server start --dev

# 停止服务器
servon server stop

# 重启服务器
servon server restart
```

### 命令参数

#### start 命令参数
- `--port, -p`: 指定服务器监听端口（默认: 9876）
- `--host`: 指定服务器监听地址（默认: 127.0.0.1）
- `--api-only`: 仅启动 API 服务，不提供 Web 界面
- `--dev`: 启用开发模式
- `--verbose, -v`: 启用详细日志模式

#### stop/restart 命令参数
- `--verbose, -v`: 启用详细日志模式

### Web 界面访问

服务器启动后，可以通过浏览器访问 Web 管理面板：

```
http://127.0.0.1:9876
```

## API 接口

### 核心 API 路由

- `/web_api/services/*` - 服务管理 API
- `/web_api/files/*` - 文件管理 API
- `/web_api/info/*` - 系统信息 API
- `/web_api/processes/*` - 进程管理 API
- `/web_api/ports/*` - 端口管理 API
- `/web_api/logs/*` - 日志管理 API
- `/web_api/soft/*` - 软件管理 API
- `/web_api/tasks/*` - 任务管理 API
- `/web_api/users/*` - 用户管理 API
- `/web_api/github/*` - GitHub 集成 API
- `/web_api/deploy/*` - 部署管理 API
- `/web_api/topology/*` - 拓扑管理 API
- `/cron/*` - 定时任务 API

### 主要 API 端点示例

```bash
# 获取系统资源信息
GET /web_api/info/resources

# 获取文件列表
GET /web_api/files/

# 获取服务列表
GET /web_api/services

# 启动服务
POST /web_api/services/{name}/start

# 获取进程列表
GET /web_api/processes/

# 获取日志文件列表
GET /web_api/logs/files
```

## 技术架构

### 后端技术栈
- **Web 框架**: Gin (Go)
- **路由管理**: 模块化路由设计
- **控制器模式**: MVC 架构
- **管理器模式**: 业务逻辑封装

### 前端技术栈
- **静态资源**: 嵌入式文件系统
- **单页应用**: SPA 架构
- **API 通信**: RESTful API

### 项目结构

```
plugins/server/
├── root.go              # 插件入口和命令注册
├── start.go             # 启动命令实现
├── stop.go              # 停止命令实现
├── restart.go           # 重启命令实现
├── dev.go               # 开发命令实现
└── web/                 # Web 相关代码
    ├── controllers/     # 控制器层
    │   ├── service_controller.go
    │   ├── file_controller.go
    │   ├── info_controller.go
    │   └── ...
    ├── routers/         # 路由层
    │   ├── root.go
    │   ├── service_router.go
    │   └── ...
    └── static/          # 静态资源
        └── dist/
```

## 开发说明

### 添加新功能

1. **创建控制器**: 在 `web/controllers/` 目录下创建新的控制器
2. **定义路由**: 在 `web/routers/` 目录下添加路由配置
3. **注册路由**: 在 `web/routers/root.go` 中注册新路由

### 控制器示例

```go
package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type ExampleController struct {
    // 依赖注入
}

func NewExampleController() *ExampleController {
    return &ExampleController{}
}

func (c *ExampleController) HandleExample(ctx *gin.Context) {
    ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
```

## 安全注意事项

- 默认绑定到 `127.0.0.1`，仅本地访问
- 生产环境使用时建议配置防火墙规则
- 文件操作具有系统权限，请谨慎使用
- 建议在受信任的网络环境中使用

## 故障排除

### 常见问题

1. **端口被占用**
   ```bash
   # 检查端口占用
   lsof -i :9876
   # 或使用其他端口
   servon server start --port 8080
   ```

2. **权限不足**
   ```bash
   # 使用 sudo 运行（谨慎使用）
   sudo servon server start
   ```

3. **服务无法停止**
   ```bash
   # 使用详细模式查看错误
   servon server stop --verbose
   ```

## 更新日志

- 支持 PID 文件管理
- 优化进程检测逻辑
- 增强错误处理和日志记录
- 支持优雅关闭和重启
- 完善的 Web API 接口

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个插件！