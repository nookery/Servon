# Servon - 服务器管理工具

Servon 是一个多功能的服务器管理工具，提供项目部署、软件安装以及可视化管理面板等功能。

当前尚处于开发阶段，功能可能不稳定，请谨慎使用。

## 功能特性

- 项目部署（`servon deploy`）
  - 支持多种项目类型的快速部署
  - 自动配置运行环境
  - 部署过程可视化

- 软件管理（`servon install`）
  - 一键安装常用服务器软件（如 Caddy、Nginx 等）
  - 自动配置和优化
  - 版本管理

- 可视化管理面板（`servon serve`）
  - 系统资源监控（CPU、内存、磁盘使用情况）
  - 网站管理（创建、配置、部署）
  - Docker 容器管理
  - 直观的 Web 操作界面

## 快速安装

### 方法 1：一键安装（推荐）

```bash
curl -fsSL https://raw.githubusercontent.com/nookery/servon/main/install.sh | bash
```

### 方法 2：手动安装

从 [GitHub Releases](https://github.com/nookery/servon/releases) 页面下载适合您系统的预编译二进制文件：

```bash
# 下载二进制文件（以 Linux amd64 为例）
curl -LO https://github.com/nookery/servon/releases/latest/download/servon-linux-amd64
chmod +x servon-linux-amd64
sudo mv servon-linux-amd64 /usr/local/bin/servon
```

## 使用方法

### 命令行界面

- 启动管理面板：`servon serve`
  - 选项：
    - `-p, --port`: 指定端口号（默认：8080）

- 查看系统信息：`servon info`

  - 选项：
    - `-f, --format`: 输出格式（formatted|json|plain）

- 实时监控：`servon monitor`

  - 选项：
    - `-i, --interval`: 监控间隔（秒）

- 查看版本：`servon version`

### Web 界面

启动服务后，访问 `http://localhost:8080` 即可使用 Web 管理界面。

## 系统要求

- 操作系统：Linux、macOS
- 建议内存：>= 512MB
- 磁盘空间：>= 200MB

## 构建说明

本项目使用 Makefile 管理构建流程。在构建之前，某些插件（如 web 插件）需要执行预处理步骤（如 npm build）。

### 构建命令

```bash 
# 完整构建（包含所有预处理步骤）
make build
# 开发模式构建（跳过预处理步骤）
make dev
# 仅执行生成步骤
make generate
# 清理构建产物
make clean
```


### 构建模式说明

1. **完整构建** (`make build`)
   - 执行所有预处理步骤（如 web 插件的 npm build）
   - 编译整个项目
   - 输出位置：`bin/servon`
   - 适用于：生产环境部署、发布新版本

2. **开发模式** (`make dev`)
   - 跳过耗时的预处理步骤
   - 直接编译项目
   - 适用于：日常开发、快速测试

3. **生成步骤** (`make generate`)
   - 仅执行预处理步骤
   - 不进行项目编译
   - 适用于：单独更新前端资源

4. **清理构建** (`make clean`)
   - 清除所有构建产物
   - 删除 `bin/` 目录
   - 删除插件生成的资源（如 `plugins/web/dist/`）

### 环境变量

- `SKIP_GENERATE`：设置此变量可跳过预处理步骤
  ```bash
  SKIP_GENERATE=1 make build  # 跳过预处理直接构建
  ```

### 插件开发者注意事项

如果你的插件需要在构建前执行特定的预处理步骤：

1. 在插件目录添加 `generate.go`：
   ```go
   //go:generate your-command
   package yourplugin
   ```

2. 创建插件级别的 Makefile：
   ```makefile
   .PHONY: build generate
   
   generate:
       @if [ "$(SKIP_GENERATE)" = "" ]; then \
           go generate ./... ; \
       fi
   
   build: generate
       go build
   ```

3. 在根目录的 Makefile 中添加你的插件：
   ```makefile
   generate:
       @echo "Generating assets..."
       @cd plugins/web && make generate
       @cd plugins/your-plugin && make generate  # 添加这行
   ```

## 贡献

欢迎提交 Pull Request 或创建 Issue！

## 许可证

MIT License
