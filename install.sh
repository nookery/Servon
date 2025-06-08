#!/bin/bash

# Servon 安装脚本
# 
# 使用方法:
#   curl -fsSL https://raw.githubusercontent.com/nookery/servon/main/install.sh | bash
#
# 环境变量:
#   GITHUB_TOKEN    - GitHub Personal Access Token，用于提高 API 速率限制
#   SERVON_VERSION  - 指定要安装的版本，例如 v1.0.0（跳过 API 调用）
#   DEBUG          - 启用调试模式，显示详细的API响应信息
#
# 示例:
#   # 使用认证安装最新版本
#   GITHUB_TOKEN=your_token bash install.sh
#
#   # 安装指定版本
#   SERVON_VERSION=v1.0.0 bash install.sh
#
#   # 在 GitHub Actions 中使用
#   env:
#     GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
#   run: curl -fsSL https://raw.githubusercontent.com/nookery/servon/main/install.sh | bash
#
#   # 启用调试模式
#   DEBUG=1 bash install.sh

# 颜色定义
RED='\033[0;31m'      # 红色文字
GREEN='\033[0;32m'    # 绿色文字
YELLOW='\033[1;33m'   # 黄色文字
BLUE='\033[0;34m'     # 蓝色文字
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    printf "%b🔍 %s%b\n" "${BLUE}" "$1" "${NC}"
}

print_success() {
    printf "%b✅ %s%b\n" "${GREEN}" "$1" "${NC}"
}

print_error() {
    printf "%b❌ %s%b\n" "${RED}" "$1" "${NC}"
}

print_warning() {
    printf "%b⚠️ %s%b\n" "${YELLOW}" "$1" "${NC}"
}

# 错误处理函数
handle_error() {
    local exit_code=$1
    local command=$2
    
    print_error "Command failed: $command"
    print_error "Exit code: $exit_code"
    
    # 如果是 curl 命令失败，显示详细信息
    if [[ "$command" == *"curl"* ]]; then
        print_error "Trying with verbose output:"
        curl -v "https://api.github.com/repos/nookery/Servon/releases/latest"
    fi
    
    exit $exit_code
}

# 执行命令并检查结果
run_command() {
    "$@"
    local exit_code=$?
    if [ $exit_code -ne 0 ]; then
        handle_error $exit_code "$*"
    fi
    return $exit_code
}

# 检查并以 sudo 权限运行命令
run_with_sudo() {
    if [ "$(id -u)" -eq 0 ]; then
        "$@"
    else
        sudo "$@"
    fi
}

# 检查系统要求
check_system() {
    print_info "Checking system requirements..."
    
    # 检查操作系统
    OS=$(uname -s)
    case "$OS" in
        Linux)
            # 检查是否为 Ubuntu
            if [ -f /etc/os-release ]; then
                . /etc/os-release
                if [ "$ID" != "ubuntu" ]; then
                    print_error "This script only supports Ubuntu on Linux. Your system ($ID) is not supported."
                    exit 1
                fi
            else
                print_error "Cannot determine Linux distribution. This script only supports Ubuntu."
                exit 1
            fi
            ;;
        Darwin)
            print_info "Detected macOS system"
            ;;
        *)
            print_error "This script only supports Linux (Ubuntu) and macOS. Your system ($OS) is not supported."
            exit 1
            ;;
    esac

    # 检查必要的命令
    for cmd in curl; do
        if ! command -v "$cmd" > /dev/null 2>&1; then
            print_error "$cmd is required but not installed"
            exit 1
        fi
    done

    # 在 macOS 上检查 shasum，在 Linux 上检查 sha256sum
    if [ "$OS" = "Darwin" ]; then
        if ! command -v "shasum" > /dev/null 2>&1; then
            print_error "shasum is required but not installed"
            exit 1
        fi
    else
        if ! command -v "sha256sum" > /dev/null 2>&1; then
            print_error "sha256sum is required but not installed"
            exit 1
        fi
    fi

    # 如果当前用户是 root，则不需要检查 sudo
    if [ "$(id -u)" -ne 0 ]; then
        if ! command -v "sudo" > /dev/null 2>&1; then
            print_error "sudo is required but not installed"
            exit 1
        fi
    fi

    print_success "System requirements met"
}

# 创建安装目录
create_install_dir() {
    print_info "Creating installation directory..."
    
    # 为不同的操作系统设置不同的安装路径
    local os=$(uname -s)
    case "$os" in
        Darwin)
            INSTALL_DIR="/usr/local/servon"
            ;;
        Linux)
            INSTALL_DIR="/usr/local/servon"
            ;;
        *)
            print_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac

    # 创建安装目录
    run_with_sudo mkdir -p "$INSTALL_DIR"

    if [ $? -ne 0 ]; then
        print_error "Failed to create installation directory"
        exit 1
    fi

    print_success "Installation directory created at $INSTALL_DIR"
}

# 检测系统架构
detect_arch() {
    local arch
    arch=$(uname -m)
    case $arch in
        x86_64)  echo "amd64" ;;
        aarch64) echo "arm64" ;;
        arm64)   echo "arm64" ;;  # macOS M1/M2 芯片
        *)       echo "unsupported" ;;
    esac
}

# 检测操作系统
detect_os() {
    local os
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$os" in
        linux)
            echo "linux"
            ;;
        darwin)
            echo "darwin"
            ;;
        *)
            print_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
}

# 获取最新版本（支持 GitHub Token 认证）
get_latest_version() {
    local api_url="https://api.github.com/repos/nookery/Servon/releases/latest"
    
    # 检查是否提供了 GitHub Token
    if [ -n "$GITHUB_TOKEN" ]; then
        # 使用认证请求
        local api_response
        api_response=$(curl -s -w "\n%{http_code}" \
            -H "Authorization: Bearer $GITHUB_TOKEN" \
            -H "User-Agent: Servon-Installer" \
            "$api_url")
    else
        # 使用未认证请求
        local api_response
        api_response=$(curl -s -w "\n%{http_code}" "$api_url")
    fi
    
    local status_code=$(echo "$api_response" | tail -n1)
    local response_body=$(echo "$api_response" | sed '$d')

    # 调试信息：显示响应的前几行
    if [ -n "$DEBUG" ]; then
        print_info "Debug: HTTP Status Code: $status_code" >&2
        print_info "Debug: Response body (first 200 chars): $(echo "$response_body" | head -c 200)..." >&2
    fi

    # 检查 HTTP 状态码
    if [ "$status_code" != "200" ]; then
        print_error "Failed to fetch latest version. HTTP Status: $status_code" >&2
        print_error "API Response: $response_body" >&2
        
        # 如果是速率限制错误，提供解决建议
        if echo "$response_body" | grep -q "rate limit exceeded"; then
            print_error "GitHub API rate limit exceeded!" >&2
            print_info "Solutions:" >&2
            print_info "1. Set GITHUB_TOKEN environment variable for authentication" >&2
            print_info "2. Wait and retry later (resets every hour)" >&2
            print_info "3. Use direct download with specific version" >&2
            print_info "   Visit https://github.com/nookery/Servon/releases to find the latest version" >&2
            print_info "   Example: SERVON_VERSION=v1.0.0 bash install.sh" >&2
        fi
        
        return 1
    fi

    # 尝试获取版本号
    local version
    # 使用更精确的JSON解析，避免匹配到错误内容
    version=$(echo "$response_body" | grep -o '"tag_name"[[:space:]]*:[[:space:]]*"[^"]*"' | sed -E 's/.*"([^"]+)"/\1/')
    
    # 如果第一种方法失败，尝试备用解析方法
    if [ -z "$version" ]; then
        version=$(echo "$response_body" | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n1)
    fi
    
    # 验证版本号格式（应该是v开头的版本号或纯数字版本号）
    if [ -n "$version" ] && ! echo "$version" | grep -qE '^v?[0-9]+\.[0-9]+\.[0-9]+'; then
        print_warning "Parsed version '$version' doesn't match expected format" >&2
        print_error "API Response: $response_body" >&2
        version=""
    fi
    
    if [ -z "$version" ]; then
        print_error "No valid version tag found in the API response" >&2
        print_error "API Response: $response_body" >&2
        return 1
    fi

    echo "$version"
}

# 下载最新版本
download_latest() {
    local version
    
    # 检查是否指定了版本
    if [ -n "$SERVON_VERSION" ]; then
        version="$SERVON_VERSION"
        print_info "Using specified version: $version"
    else
        print_info "Fetching latest version from GitHub API..."
        version=$(get_latest_version)
        if [ -z "$version" ]; then
            print_error "Failed to get latest version"
            exit 1
        fi
        print_success "Found latest version: $version"
    fi

    local arch os
    arch=$(detect_arch)
    os=$(detect_os)
    
    # 构建下载 URL
    local download_url="https://github.com/nookery/Servon/releases/download/${version}/servon-${os}-${arch}"
    local checksum_url="https://github.com/nookery/Servon/releases/download/${version}/servon-${os}-${arch}.sha256"
    
    print_info "Downloading Servon ${version} for ${os}/${arch}..."

    # 下载二进制文件
    print_info "Downloading binary..."
    curl -L -o "/tmp/servon-${os}-${arch}" "$download_url"
    if [ $? -ne 0 ]; then
        print_error "Failed to download binary from: $download_url"
        exit 1
    fi

    # 下载校验和文件
    print_info "Downloading checksum..."
    curl -L -o "/tmp/servon-${os}-${arch}.sha256" "$checksum_url"
    if [ $? -ne 0 ]; then
        print_error "Failed to download checksum from: $checksum_url"
        exit 1
    fi

    # 验证校验和
    print_info "Verifying download..."
    if [ "$os" = "darwin" ]; then
        # macOS 使用 shasum
        if ! (cd /tmp && shasum -a 256 -c "servon-${os}-${arch}.sha256"); then
            print_error "Checksum verification failed"
            exit 1
        fi
    else
        # Linux 使用 sha256sum
        if ! (cd /tmp && sha256sum -c "servon-${os}-${arch}.sha256"); then
            print_error "Checksum verification failed"
            exit 1
        fi
    fi

    print_success "Downloaded and verified Servon ${version}"
}

# 安装文件
install_files() {
    print_info "Installing Servon..."
    
    local arch=$(detect_arch)
    local os=$(detect_os)
    
    # 移动二进制文件
    run_with_sudo mv "/tmp/servon-${os}-${arch}" "$INSTALL_DIR/servon"
    if [ $? -ne 0 ]; then
        print_error "Failed to install Servon"
        exit 1
    fi

    # 设置执行权限
    run_with_sudo chmod +x "$INSTALL_DIR/servon"
    if [ $? -ne 0 ]; then
        print_error "Failed to set permissions"
        exit 1
    fi

    # 创建符号链接
    run_with_sudo ln -sf "$INSTALL_DIR/servon" "/usr/local/bin/servon"
    if [ $? -ne 0 ]; then
        print_error "Failed to create symbolic link"
        exit 1
    fi

    # 清理临时文件
    rm -f "/tmp/servon-${os}-${arch}.sha256"

    print_success "Servon installed successfully"
}

# 打印使用说明
print_usage() {
    echo
    print_info "Servon has been installed successfully!"
    echo
    echo "You can now use Servon with the following commands:"
    echo
    echo "  servon"
    echo
    echo "For more information, visit: https://github.com/nookery/servon"
    echo
}

# 主函数
main() {
    echo "=============================="
    echo "  Servon Installer"
    echo "=============================="
    echo

    check_system
    create_install_dir
    download_latest
    install_files
    print_usage
}

# 运行主函数
main "$@"