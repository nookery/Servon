# GitHub API 速率限制解决方案

## 问题描述

在 GitHub Actions 中使用 `install.sh` 脚本时，遇到 GitHub API 速率限制错误：

```
❌ API Response: {"message":"API rate limit exceeded for 13.105.117.149. (But here's the good news: Authenticated requests get a higher rate limit. Check out the documentation for more details.)","documentation_url":"https://docs.github.com/rest/overview/resources-in-the-rest-api#rate-limiting"}
```

## 根本原因

1. **未认证请求限制**：未认证的 GitHub API 请求每小时限制为 60 次
2. **GitHub Actions IP 共享**：多个 GitHub Actions 实例可能共享相同的出口 IP
3. **频繁调用**：多个项目同时使用相同的安装脚本

## 解决方案

### 方案一：指定版本安装（最简单，推荐）

通过设置 `SERVON_VERSION` 环境变量直接指定要安装的版本，完全跳过 GitHub API 调用：

```yaml
- name: 安装 Servon
  env:
    SERVON_VERSION: v1.0.0  # 替换为实际版本
  run: |
    curl -fsSL 'https://raw.githubusercontent.com/nookery/servon/main/install.sh' | bash
```

**优点：**
- 完全避免 API 速率限制问题
- 安装速度更快（无需 API 调用）
- 版本固定，构建更稳定
- 无需额外配置

**缺点：**
- 需要手动更新版本号
- 无法自动获取最新版本

### 方案二：使用 GitHub Token 认证（动态版本）

#### 1. 修改 install.sh 脚本

在 `get_latest_version()` 函数中添加认证支持：

```bash
# 获取最新版本（支持认证）
get_latest_version() {
    local api_url="https://api.github.com/repos/nookery/Servon/releases/latest"
    local curl_args="-s -w \"\n%{http_code}\""
    
    # 如果提供了 GitHub Token，则使用认证
    if [ -n "$GITHUB_TOKEN" ]; then
        curl_args="$curl_args -H \"Authorization: Bearer $GITHUB_TOKEN\""
        curl_args="$curl_args -H \"User-Agent: Servon-Installer\""
        print_info "Using authenticated GitHub API request"
    else
        print_warning "Using unauthenticated GitHub API request (rate limited)"
    fi
    
    # 执行 API 请求
    local api_response
    if [ -n "$GITHUB_TOKEN" ]; then
        api_response=$(curl -s -w "\n%{http_code}" \
            -H "Authorization: Bearer $GITHUB_TOKEN" \
            -H "User-Agent: Servon-Installer" \
            "$api_url")
    else
        api_response=$(curl -s -w "\n%{http_code}" "$api_url")
    fi
    
    local status_code=$(echo "$api_response" | tail -n1)
    local response_body=$(echo "$api_response" | sed '$d')

    # 检查 HTTP 状态码
    if [ "$status_code" != "200" ]; then
        print_error "Failed to fetch latest version. HTTP Status: $status_code"
        print_error "API Response: $response_body"
        
        # 如果是速率限制错误，提供解决建议
        if echo "$response_body" | grep -q "rate limit exceeded"; then
            print_error "GitHub API rate limit exceeded!"
            print_info "Solutions:"
            print_info "1. Set GITHUB_TOKEN environment variable"
            print_info "2. Wait and retry later"
            print_info "3. Use direct download with specific version"
        fi
        
        return 1
    fi

    # 尝试获取版本号
    local version
    version=$(echo "$response_body" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$version" ]; then
        print_error "No version tag found in the API response"
        print_error "API Response: $response_body"
        return 1
    fi

    echo "$version"
}
```

#### 2. 在 GitHub Actions 中使用

```yaml
- name: 安装依赖
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  run: |
    curl -fsSL \
      --header 'authorization: Bearer ${{ secrets.GITHUB_TOKEN }}' \
      --header 'user-agent: GitOK-Release-Workflow' \
      'https://raw.githubusercontent.com/nookery/servon/main/install.sh' | bash
```

### 方案三：使用缓存机制

#### 在 GitHub Actions 中缓存下载的文件

```yaml
- name: 缓存 Servon 二进制文件
  id: cache-servon
  uses: actions/cache@v3
  with:
    path: ~/.local/bin/servon
    key: servon-${{ runner.os }}-${{ runner.arch }}-latest
    restore-keys: |
      servon-${{ runner.os }}-${{ runner.arch }}-

- name: 安装 Servon
  if: steps.cache-servon.outputs.cache-hit != 'true'
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  run: |
    curl -fsSL \
      --header 'authorization: Bearer ${{ secrets.GITHUB_TOKEN }}' \
      --header 'user-agent: GitOK-Release-Workflow' \
      'https://raw.githubusercontent.com/nookery/servon/main/install.sh' | bash
```

### 方案四：直接下载特定版本

#### 避免 API 调用，直接下载已知版本

```yaml
- name: 安装 Servon (特定版本)
  run: |
    VERSION="v1.0.0"  # 指定版本
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    case $ARCH in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
    esac
    
    # 直接下载二进制文件
    curl -L -o servon \
      "https://github.com/nookery/servon/releases/download/${VERSION}/servon-${OS}-${ARCH}"
    
    chmod +x servon
    sudo mv servon /usr/local/bin/
```

### 方案五：使用 GitHub CLI

#### 利用 GitHub CLI 的内置认证

```yaml
- name: 安装 Servon 使用 GitHub CLI
  run: |
    # 获取最新版本
    VERSION=$(gh release view --repo nookery/servon --json tagName --jq '.tagName')
    
    # 下载二进制文件
    gh release download $VERSION --repo nookery/servon --pattern "servon-*" --dir /tmp
    
    # 安装
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    case $ARCH in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
    esac
    
    sudo mv "/tmp/servon-${OS}-${ARCH}" /usr/local/bin/servon
    sudo chmod +x /usr/local/bin/servon
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### 方案六：创建专用的 GitHub Action

#### 创建可重用的 Action

创建 `.github/actions/install-servon/action.yml`：

```yaml
name: 'Install Servon'
description: 'Install Servon with rate limit handling'
inputs:
  version:
    description: 'Servon version to install'
    required: false
    default: 'latest'
  github-token:
    description: 'GitHub token for API access'
    required: false
    default: ${{ github.token }}
runs:
  using: 'composite'
  steps:
    - name: Cache Servon
      id: cache
      uses: actions/cache@v3
      with:
        path: ~/.local/bin/servon
        key: servon-${{ runner.os }}-${{ runner.arch }}-${{ inputs.version }}
    
    - name: Install Servon
      if: steps.cache.outputs.cache-hit != 'true'
      shell: bash
      env:
        GITHUB_TOKEN: ${{ inputs.github-token }}
      run: |
        if [ "${{ inputs.version }}" = "latest" ]; then
          # 使用认证的 API 请求
          VERSION=$(curl -s -H "Authorization: Bearer $GITHUB_TOKEN" \
            https://api.github.com/repos/nookery/servon/releases/latest | \
            grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
        else
          VERSION="${{ inputs.version }}"
        fi
        
        OS=$(uname -s | tr '[:upper:]' '[:lower:]')
        ARCH=$(uname -m)
        case $ARCH in
            x86_64) ARCH="amd64" ;;
            aarch64|arm64) ARCH="arm64" ;;
        esac
        
        # 下载并安装
        curl -L -o ~/.local/bin/servon \
          "https://github.com/nookery/servon/releases/download/${VERSION}/servon-${OS}-${ARCH}"
        chmod +x ~/.local/bin/servon
        echo "$HOME/.local/bin" >> $GITHUB_PATH
```

使用方式：

```yaml
- name: Install Servon
  uses: ./.github/actions/install-servon
  with:
    version: 'latest'
    github-token: ${{ secrets.GITHUB_TOKEN }}
```

## 推荐方案总结

根据不同场景，推荐使用以下方案：

### 🥇 生产环境（推荐）
**方案一：指定版本安装**
- 完全避免 API 限制
- 构建稳定可重现
- 无需额外配置

### 🥈 开发环境
**方案二：GitHub Token 认证**
- 自动获取最新版本
- 高速率限制（5000/小时）
- 适合频繁更新

### 🥉 备选方案
**方案三：缓存机制**
- 减少重复下载
- 提高构建速度
- 适合大型项目

## 最佳实践建议

1. **生产优先固定版本**：使用 `SERVON_VERSION` 环境变量
2. **开发使用认证**：在 GitHub Actions 中始终使用 `GITHUB_TOKEN`
3. **实现缓存**：对于频繁构建的项目，使用缓存机制
4. **错误处理**：添加重试逻辑和降级方案
5. **监控使用**：定期检查 API 使用情况

## 常见问题

### Q: 为什么会遇到速率限制？
A: GitHub API 对未认证请求有严格限制（每小时 60 次），GitHub Actions 的多个实例可能共享 IP 地址。

### Q: 使用 GITHUB_TOKEN 安全吗？
A: 是的，`GITHUB_TOKEN` 是 GitHub Actions 自动提供的临时令牌，权限范围有限且会自动过期。

### Q: 如何获取个人访问令牌？
A: 在 GitHub 设置中创建 Personal Access Token，只需要 `public_repo` 权限即可读取公开仓库的发布信息。

### Q: SERVON_VERSION 支持哪些版本格式？
A: 支持 GitHub Release 的标签格式，如 `v1.0.0`、`1.0.0`、`latest` 等。

### 2. 组合使用

可以组合多个方案以获得最佳效果：

```yaml
- name: 安装 Servon（组合方案）
  run: |
    # 首先尝试从缓存获取
    if [ -f ~/.local/bin/servon ]; then
      echo "Using cached Servon"
      echo "$HOME/.local/bin" >> $GITHUB_PATH
      exit 0
    fi
    
    # 使用认证的安装脚本
    export GITHUB_TOKEN="${{ secrets.GITHUB_TOKEN }}"
    curl -fsSL \
      --header 'authorization: Bearer ${{ secrets.GITHUB_TOKEN }}' \
      --header 'user-agent: GitOK-Release-Workflow' \
      'https://raw.githubusercontent.com/nookery/servon/main/install.sh' | bash
    
    # 缓存到用户目录
    mkdir -p ~/.local/bin
    cp /usr/local/bin/servon ~/.local/bin/
    echo "$HOME/.local/bin" >> $GITHUB_PATH
```

### 3. 错误处理

添加重试机制和降级方案：

```yaml
- name: 安装 Servon（带重试）
  run: |
    for i in {1..3}; do
      if GITHUB_TOKEN="${{ secrets.GITHUB_TOKEN }}" \
         curl -fsSL 'https://raw.githubusercontent.com/nookery/servon/main/install.sh' | bash; then
        echo "Servon installed successfully"
        break
      else
        echo "Attempt $i failed, retrying in 30 seconds..."
        sleep 30
      fi
      
      if [ $i -eq 3 ]; then
        echo "All attempts failed, using fallback method"
        # 降级到直接下载特定版本
        VERSION="v1.0.0"
        OS=$(uname -s | tr '[:upper:]' '[:lower:]')
        ARCH=$(uname -m)
        case $ARCH in
            x86_64) ARCH="amd64" ;;
            aarch64|arm64) ARCH="arm64" ;;
        esac
        curl -L -o servon \
          "https://github.com/nookery/servon/releases/download/${VERSION}/servon-${OS}-${ARCH}"
        chmod +x servon
        sudo mv servon /usr/local/bin/
      fi
    done
```

## 总结

推荐使用 **方案一**（GitHub Token 认证）作为主要解决方案，因为：

1. 简单易实现
2. 解决根本问题
3. 提高速率限制（每小时 5000 次）
4. 不需要修改现有工作流太多

同时可以结合缓存机制进一步优化性能和可靠性。