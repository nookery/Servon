# 更新机制

## 概述

Servon采用基于GitHub Releases的版本管理和更新机制。系统通过版本管理器（VersionManager）来处理版本信息的获取、比较和更新检查。

## 版本信息

### 版本组成

Servon的版本信息包含三个主要部分：

1. **Version** - 版本号（如"1.0.0"）
2. **CommitHash** - Git提交哈希值
3. **BuildTime** - 构建时间

这些信息在构建时通过ldflags注入：

```bash
go build -ldflags "-X main.Version=1.0.0 -X main.CommitHash=abc123 -X main.BuildTime=2024-01-01"
```

### 开发版本标识

在开发环境中，系统会：
- 使用"dev"作为默认版本号
- 尝试从package.json读取版本信息
- 在版本号后添加"(dev)"标识

## 版本检查

### 获取最新版本

VersionManager提供了从GitHub获取最新发布版本的功能：

```go
latestVersion, err := versionManager.GetLatestVersion()
```

这个功能会：
1. 访问GitHub API获取最新release信息
2. 解析release的tag名称
3. 返回不带"v"前缀的版本号

### 版本信息查询

系统提供了以下方法来获取版本信息：

1. **获取当前版本号**
```go
version := versionManager.GetVersion()
```

2. **获取完整版本信息**
```go
versionInfo := versionManager.GetVersionInfo()
fmt.Printf("Version: %s\nCommit: %s\nBuild Time: %s\n",
    versionInfo.Version,
    versionInfo.CommitHash,
    versionInfo.BuildTime)
```

## 构建与发布

### 版本注入

在构建过程中，通过Makefile的LDFLAGS参数注入版本信息：

```makefile
# 构建后端
build-backend:
    go build -ldflags "$(LDFLAGS)" -o temp/servon
```

使用示例：
```bash
make build-backend LDFLAGS="-X main.Version=1.0.0 -X main.CommitHash=$(git rev-parse HEAD) -X main.BuildTime=$(date -u '+%Y-%m-%d')"
```

### 发布流程

1. **版本标记**
   - 在GitHub上创建新的release
   - 使用"v"前缀的语义化版本号（如"v1.0.0"）
   - 提供详细的更新说明

2. **构建发布**
   - 使用正确的版本信息构建
   - 将构建产物上传到GitHub release

## 最佳实践

1. **版本号管理**
   - 遵循语义化版本（Semantic Versioning）规范
   - 主版本号：不兼容的API修改
   - 次版本号：向下兼容的功能性新增
   - 修订号：向下兼容的问题修正

2. **更新检查**
   - 定期检查新版本
   - 在关键操作前验证版本兼容性
   - 提供更新提示和说明

3. **开发环境**
   - 使用dev标识区分开发版本
   - 保持package.json版本信息同步
   - 在测试环境验证版本信息正确性

## 常见问题

1. **版本检查失败**
   - 检查网络连接
   - 验证GitHub API访问权限
   - 确认release tag格式正确

2. **版本信息不正确**
   - 检查构建命令中的LDFLAGS
   - 确认版本信息正确注入
   - 验证package.json中的版本号

3. **开发环境问题**
   - 确认是否正确识别为开发环境
   - 检查package.json文件存在且可访问
   - 验证版本号格式正确