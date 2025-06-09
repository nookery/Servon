# User 用户管理插件

## 概述

User插件是Servon的用户管理工具，提供了完整的系统用户管理功能，包括用户列表查看、用户创建、用户删除和用户详细信息查看等功能。该插件支持跨平台操作，兼容macOS和Linux系统。

## 功能特性

### 🔍 用户列表查看
- 列出系统中的所有用户
- 支持过滤系统用户
- 提供详细信息显示模式
- 显示用户权限状态（sudo权限）

### 👤 用户详细信息
- 查看用户基本信息（用户名、主目录、Shell等）
- 显示用户权限信息
- 统计用户文件和目录信息
- 查看用户当前运行的进程
- 显示用户登录历史

### ➕ 用户创建
- 创建新的系统用户
- 支持密码设置（交互式输入）
- 可配置用户Shell
- 支持添加到指定用户组
- 可授予sudo权限
- 提供用户名格式验证

### 🗑️ 用户删除
- 安全删除系统用户
- 可选择是否删除用户主目录
- 提供受保护用户检查
- 删除前确认机制
- 检查用户运行进程

## 使用方法

### 基本命令结构
```bash
servon user [command] [flags]
```

### 1. 列出用户

#### 基本用法
```bash
# 列出所有非系统用户
servon user list

# 包含系统用户
servon user list --system

# 显示详细信息
servon user list --verbose

# 组合使用
servon user list -sv
```

#### 参数说明
- `-s, --system`: 包含系统用户
- `-v, --verbose`: 显示详细信息

### 2. 查看用户详细信息

#### 基本用法
```bash
# 查看用户基本信息
servon user info [username]

# 显示用户当前进程
servon user info [username] --processes

# 显示登录历史
servon user info [username] --login-history

# 显示所有信息
servon user info [username] -pl
```

#### 参数说明
- `-p, --processes`: 显示用户当前运行的进程
- `-l, --login-history`: 显示用户登录历史

#### 示例
```bash
# 查看当前用户信息
servon user info $(whoami)

# 查看用户详细信息和进程
servon user info colorfy --processes
```

### 3. 创建用户

#### 基本用法
```bash
# 创建基本用户（会提示输入密码）
servon user create [username]

# 指定密码创建用户
servon user create [username] --password "your_password"

# 创建用户并设置Shell
servon user create [username] --shell /bin/bash

# 创建用户并添加到用户组
servon user create [username] --groups docker,www-data

# 创建用户并授予sudo权限
servon user create [username] --sudo

# 强制创建（覆盖已存在用户）
servon user create [username] --force
```

#### 参数说明
- `-p, --password`: 指定用户密码
- `-s, --shell`: 指定用户Shell
- `-g, --groups`: 添加到的用户组列表（逗号分隔）
- `-S, --sudo`: 添加sudo权限
- `-f, --force`: 强制创建（覆盖已存在用户）

#### 示例
```bash
# 创建开发用户
servon user create developer --shell /bin/bash --groups docker,sudo --sudo

# 创建Web服务用户
servon user create webuser --groups www-data,nginx
```

### 4. 删除用户

#### 基本用法
```bash
# 删除用户（保留主目录）
servon user delete [username]

# 删除用户及主目录
servon user delete [username] --remove-home

# 强制删除（跳过确认）
servon user delete [username] --force

# 详细模式删除
servon user delete [username] --verbose
```

#### 参数说明
- `-r, --remove-home`: 同时删除用户主目录
- `-f, --force`: 强制删除（跳过确认和保护检查）
- `-v, --verbose`: 显示详细信息

#### 示例
```bash
# 安全删除用户
servon user delete testuser --verbose

# 完全删除用户
servon user delete olduser --remove-home --force
```

## 系统兼容性

### macOS支持
- 使用`dscl`命令获取用户信息
- 支持macOS用户目录结构
- 兼容macOS权限系统

### Linux支持
- 使用`/etc/passwd`文件读取用户信息
- 支持标准Linux用户管理命令
- 兼容各种Linux发行版

## 安全特性

### 用户创建安全
- 用户名格式验证
- 密码强度要求
- 重复密码确认
- 用户存在性检查

### 用户删除保护
- 受保护系统用户列表
- 删除前确认机制
- 进程检查警告
- 强制删除选项

### 权限管理
- sudo权限检测
- 用户组管理
- 文件权限显示

## 错误处理

插件提供了完善的错误处理机制：

- **权限不足**: 提示需要管理员权限
- **用户不存在**: 清晰的错误提示
- **命令执行失败**: 详细的错误信息
- **系统兼容性**: 自动检测操作系统类型

## 使用示例

### 日常用户管理
```bash
# 查看所有用户
servon user list -v

# 创建新的开发用户
servon user create devuser --shell /bin/bash --sudo

# 查看用户详细信息
servon user info devuser -p

# 删除测试用户
servon user delete testuser --remove-home
```

### 系统管理场景
```bash
# 查看系统用户
servon user list --system

# 创建服务用户
servon user create nginx --groups www-data

# 检查用户进程
servon user info nginx --processes

# 清理无用用户
servon user delete oldservice --force --verbose
```

## 注意事项

1. **权限要求**: 用户创建和删除操作需要管理员权限
2. **系统用户**: 默认过滤系统用户，使用`--system`参数显示
3. **密码安全**: 建议使用交互式密码输入而非命令行参数
4. **备份重要**: 删除用户前请确保重要数据已备份
5. **进程检查**: 删除用户前检查是否有运行的进程

## 故障排除

### 常见问题

**Q: 无法获取用户列表**
A: 检查系统权限，确保有读取用户信息的权限

**Q: 创建用户失败**
A: 确保有管理员权限，检查用户名格式是否正确

**Q: 删除用户时提示受保护**
A: 系统用户受保护，使用`--force`参数强制删除（谨慎操作）

**Q: macOS上显示用户信息不完整**
A: 某些用户信息可能需要特殊权限，这是正常现象

## 更新日志

### v1.0.0
- 初始版本发布
- 支持基本用户管理功能
- 兼容macOS和Linux系统
- 提供完整的命令行接口

## 贡献

欢迎提交Issue和Pull Request来改进这个插件。在提交代码前，请确保：

1. 代码符合项目规范
2. 添加适当的测试
3. 更新相关文档
4. 测试跨平台兼容性