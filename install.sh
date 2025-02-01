#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# 检查系统要求
check_system() {
    print_info "Checking system requirements..."
    
    # 检查操作系统
    if [[ "$OSTYPE" != "linux-gnu"* ]] && [[ "$OSTYPE" != "darwin"* ]]; then
        print_error "This script only supports Linux and macOS"
        exit 1
    fi

    # 检查必要的命令
    for cmd in curl sudo; do
        if ! command -v $cmd &> /dev/null; then
            print_error "$cmd is required but not installed"
            exit 1
        fi
    done

    print_success "System requirements met"
}

# 创建安装目录
create_install_dir() {
    print_info "Creating installation directory..."
    
    # 为不同的操作系统设置不同的安装路径
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        INSTALL_DIR="/usr/local/servon"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        INSTALL_DIR="/usr/local/servon"
    fi

    # 创建安装目录
    sudo mkdir -p "$INSTALL_DIR"
    if [ $? -ne 0 ]; then
        print_error "Failed to create installation directory"
        exit 1
    fi

    print_success "Installation directory created at $INSTALL_DIR"
}

# 下载最新版本
download_latest() {
    print_info "Downloading latest version..."
    
    # 获取最新版本
    LATEST_VERSION=$(curl -s https://api.github.com/repos/angel/servon/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$LATEST_VERSION" ]; then
        print_error "Failed to get latest version"
        exit 1
    fi

    # 根据系统架构选择下载文件
    ARCH=$(uname -m)
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    
    if [ "$ARCH" = "x86_64" ]; then
        ARCH="amd64"
    elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
        ARCH="arm64"
    fi

    DOWNLOAD_URL="https://github.com/angel/servon/releases/download/${LATEST_VERSION}/servon-${OS}-${ARCH}"
    CHECKSUM_URL="https://github.com/angel/servon/releases/download/${LATEST_VERSION}/servon-${OS}-${ARCH}.sha256"
    
    print_info "Downloading Servon ${LATEST_VERSION} for ${OS}/${ARCH}..."

    # 下载二进制文件
    curl -L -o "/tmp/servon" "$DOWNLOAD_URL"
    if [ $? -ne 0 ]; then
        print_error "Failed to download Servon"
        exit 1
    fi

    # 下载校验和文件
    curl -L -o "/tmp/servon.sha256" "$CHECKSUM_URL"
    if [ $? -ne 0 ]; then
        print_error "Failed to download checksum file"
        exit 1
    }

    # 验证校验和
    print_info "Verifying download..."
    if ! (cd /tmp && sha256sum -c servon.sha256); then
        print_error "Checksum verification failed"
        exit 1
    fi

    print_success "Downloaded and verified Servon ${LATEST_VERSION}"
}

# 安装文件
install_files() {
    print_info "Installing Servon..."
    
    # 移动二进制文件
    sudo mv "/tmp/servon" "$INSTALL_DIR/servon"
    if [ $? -ne 0 ]; then
        print_error "Failed to install Servon"
        exit 1
    fi

    # 设置执行权限
    sudo chmod +x "$INSTALL_DIR/servon"
    if [ $? -ne 0 ]; then
        print_error "Failed to set permissions"
        exit 1
    }

    # 创建符号链接
    sudo ln -sf "$INSTALL_DIR/servon" "/usr/local/bin/servon"
    if [ $? -ne 0 ]; then
        print_error "Failed to create symbolic link"
        exit 1
    fi

    # 清理临时文件
    rm -f "/tmp/servon.sha256"

    print_success "Servon installed successfully"
}

# 打印使用说明
print_usage() {
    echo
    print_info "Servon has been installed successfully!"
    echo
    echo "You can now use Servon with the following commands:"
    echo
    echo "  servon serve      Start the web server"
    echo "  servon info       Display system information"
    echo "  servon monitor    Monitor system resources"
    echo "  servon version    Show version information"
    echo
    echo "For more information, visit: https://github.com/angel/servon"
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
main 