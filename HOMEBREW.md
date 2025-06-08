# Homebrew 发布指南

本指南将帮助您将 Servon 项目发布到 Homebrew，让用户可以通过 `brew install servon` 来安装。

## 方式一：创建自己的 Tap（推荐）

### 1. 创建 Homebrew Tap 仓库

在 GitHub 上创建一个新的仓库，命名为 `homebrew-servon`（必须以 `homebrew-` 开头）。

### 2. 初始化 Tap

```bash
# 创建本地 tap
brew tap-new your-username/homebrew-servon

# 进入 tap 目录
cd $(brew --repository)/Library/Taps/your-username/homebrew-servon

# 添加远程仓库
git remote add origin https://github.com/your-username/homebrew-servon.git
```

### 3. 复制 Formula 文件

将项目根目录下的 `servon.rb` 文件复制到 tap 的 `Formula` 目录：

```bash
cp /path/to/servon/servon.rb Formula/servon.rb
```

### 4. 更新 Formula

编辑 `Formula/servon.rb` 文件，更新以下内容：

- 将 `homepage` 更改为您的实际 GitHub 仓库地址
- 将 `url` 更改为您的发布版本下载地址
- 计算并填入正确的 `sha256` 值
- 根据需要调整 `license`

### 5. 计算 SHA256

```bash
# 下载您的发布包
wget https://github.com/your-username/servon/archive/refs/tags/v1.0.0.tar.gz

# 计算 SHA256
shasum -a 256 v1.0.0.tar.gz
```

### 6. 测试 Formula

```bash
# 审核 formula
brew audit --strict --online servon

# 测试安装
brew install --build-from-source servon

# 测试功能
brew test servon

# 卸载测试
brew uninstall servon
```

### 7. 发布 Tap

```bash
# 提交更改
git add Formula/servon.rb
git commit -m "Add servon formula"
git push origin main
```

### 8. 用户安装方式

用户现在可以通过以下方式安装：

```bash
# 添加您的 tap
brew tap your-username/servon

# 安装 servon
brew install servon
```

## 方式二：提交到 Homebrew Core（官方）

### 前提条件

要提交到 Homebrew Core，您的项目需要满足以下条件：

1. 项目必须是开源的
2. 项目必须有稳定的发布版本
3. 项目必须有一定的知名度和用户基础
4. 项目不能是重复的工具

### 提交流程

1. **Fork homebrew-core 仓库**
   ```bash
   # Fork https://github.com/Homebrew/homebrew-core
   ```

2. **克隆并设置**
   ```bash
   # 设置环境变量
   export HOMEBREW_NO_INSTALL_FROM_API=1
   
   # Tap homebrew-core
   brew tap --force homebrew/core
   
   # 进入 homebrew-core 目录
   cd $(brew --repository homebrew/core)
   ```

3. **创建新分支**
   ```bash
   git checkout -b servon-formula origin/master
   ```

4. **添加 Formula**
   ```bash
   # 复制 formula 文件
   cp /path/to/servon/servon.rb Formula/servon.rb
   ```

5. **测试和审核**
   ```bash
   # 审核
   brew audit --strict --new --online servon
   
   # 测试安装
   brew install --build-from-source servon
   
   # 测试功能
   brew test servon
   ```

6. **提交 Pull Request**
   ```bash
   git add Formula/servon.rb
   git commit -m "servon: new formula
   
   A powerful server management and development tool
   
   Closes #xxxxx"
   git push origin servon-formula
   ```

7. **创建 Pull Request**
   - 访问 https://github.com/Homebrew/homebrew-core
   - 创建从您的分支到 master 的 Pull Request
   - 填写详细的描述信息

## Formula 文件说明

### 基本结构

```ruby
class Servon < Formula
  desc "项目描述"
  homepage "项目主页"
  url "下载地址"
  sha256 "文件校验和"
  license "许可证"
  
  # 构建依赖
  depends_on "go" => :build
  depends_on "node" => :build
  depends_on "pnpm" => :build
  
  def install
    # 安装逻辑
  end
  
  def test
    # 测试逻辑
  end
end
```

### 重要字段说明

- `desc`: 简短的项目描述
- `homepage`: 项目主页 URL
- `url`: 源码下载地址（通常是 GitHub release）
- `sha256`: 下载文件的 SHA256 校验和
- `license`: 项目许可证（使用 SPDX 标识符）
- `depends_on`: 依赖项声明

## 自动化更新

### 使用 GitHub Actions

可以创建 GitHub Actions 来自动更新 Homebrew Formula：

```yaml
name: Update Homebrew Formula

on:
  release:
    types: [published]

jobs:
  homebrew:
    runs-on: ubuntu-latest
    steps:
      - uses: mislav/bump-homebrew-formula-action@v2
        with:
          formula-name: servon
          homebrew-tap: your-username/homebrew-servon
          base-branch: main
          download-url: https://github.com/your-username/servon/archive/refs/tags/${{ github.ref_name }}.tar.gz
        env:
          COMMITTER_TOKEN: ${{ secrets.COMMITTER_TOKEN }}
```

## 维护和更新

### 版本更新

当发布新版本时：

1. 更新 `url` 中的版本号
2. 重新计算并更新 `sha256`
3. 测试新版本的安装
4. 提交更改

### 常见问题

1. **构建失败**
   - 检查依赖项是否正确
   - 确认构建命令是否正确
   - 查看构建日志定位问题

2. **测试失败**
   - 确保测试命令能正常执行
   - 检查二进制文件是否正确安装

3. **审核失败**
   - 按照 `brew audit` 的提示修复问题
   - 确保遵循 Homebrew 的编码规范

## 参考资源

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Adding Software to Homebrew](https://docs.brew.sh/Adding-Software-to-Homebrew)
- [Homebrew 官方文档](https://docs.brew.sh/)
- [Formula API 文档](https://rubydoc.brew.sh/Formula)