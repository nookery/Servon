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

本项目使用 Makefile 管理构建流程，查看 [Makefile](Makefile) 了解更多。

## 贡献

欢迎提交 Pull Request 或创建 Issue！

## 许可证

MIT License
