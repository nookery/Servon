# Servon - 服务器管理面板

Servon 是一个轻量级的服务器管理面板，提供了直观的 Web 界面来管理您的服务器。

## 功能特性

- 系统资源监控（CPU、内存、磁盘使用情况）
- 网站管理（创建、配置、部署）
- Docker 容器管理
- 可视化的 Web 界面
- 命令行工具

## 快速安装

### 方法 1：一键安装（推荐）

```bash
curl -fsSL https://raw.githubusercontent.com/angel/servon/main/install.sh | bash
```

### 方法 2：手动安装

从 [GitHub Releases](https://github.com/angel/servon/releases) 页面下载适合您系统的预编译二进制文件：

```bash
# 下载二进制文件（以 Linux amd64 为例）
curl -LO https://github.com/angel/servon/releases/latest/download/servon-linux-amd64
chmod +x servon-linux-amd64
sudo mv servon-linux-amd64 /usr/local/bin/servon
```

## 使用方法

### 命令行界面

- 启动服务：`servon serve`

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

## 开发指南

如果您想参与开发，需要安装以下依赖：

### 后端
- Go >= 1.21
- Git

### 前端
- Node.js >= 16
- pnpm >= 8.0

1. 克隆仓库：

   ```bash
   git clone https://github.com/angel/servon.git
   cd servon
   ```

2. 安装后端依赖：

   ```bash
   go mod download
   ```

3. 安装前端依赖：

   ```bash
   cd web
   pnpm install
   ```

4. 启动开发服务器：

   ```bash
   # 后端开发服务器
   go run main.go serve

   # 前端开发服务器
   cd web
   pnpm dev
   ```

5. 构建项目：
   ```bash
   # 构建后端
   go build -o servon

   # 构建前端
   cd web
   pnpm build
   ```

## 发布流程

1. 更新版本号
2. 创建 tag 并推送
3. GitHub Actions 会自动：
   - 运行测试
   - 构建多平台二进制文件
   - 生成 SHA256 校验和
   - 创建 GitHub Release
   - 上传构建产物

## 贡献

欢迎提交 Pull Request 或创建 Issue！

## 许可证

MIT License
