#!/bin/bash

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 输出步骤信息
step() {
    printf "${BLUE}[Step $1/${TOTAL_STEPS}]${NC} $2\n"
}

# 输出成功信息
success() {
    printf "${GREEN}✓${NC} $1\n"
}

# 输出错误信息
error() {
    printf "${RED}✗ Error:${NC} $1\n"
    exit 1
}

# 检查命令是否可用
check_command() {
    if command -v $1 >/dev/null 2>&1 || type $1 >/dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# 检查 pnpm 是否可用
check_pnpm() {
    if pnpm -v >/dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# 总步骤数
TOTAL_STEPS=6

# 步骤 1: 检查系统要求
step 1 "检查系统要求..."

if ! check_command node; then
    printf "  • 未检测到 Node.js，正在安装...\n"
    curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
    sudo apt-get install -y nodejs
    success "Node.js 安装完成"
else
    success "Node.js 已安装"
fi

if ! check_pnpm; then
    printf "  • 未检测到 pnpm，正在安装...\n"
    curl -fsSL https://get.pnpm.io/install.sh | sh -
    success "pnpm 安装完成"
else
    success "pnpm 已安装"
fi

if ! command -v systemctl &> /dev/null; then
    error "未检测到 systemd。本程序需要 systemd 支持，请确保你的系统支持 systemd"
fi

# 步骤 2: 下载项目
step 2 "下载项目代码..."
git clone https://github.com/nookery/servon.git
if [ $? -ne 0 ]; then
    error "无法下载项目代码。请检查你的网络连接或确认 GitHub 是否可访问"
fi
success "项目代码下载完成"

cd servon
if [ $? -ne 0 ]; then
    error "无法进入项目目录。请检查目录权限"
fi

# 步骤 3: 安装依赖
step 3 "安装项目依赖..."
pnpm install
if [ $? -ne 0 ]; then
    error "依赖安装失败。请检查你的网络连接或尝试清除 pnpm 缓存后重试"
fi
success "依赖安装完成"

# 步骤 4: 构建项目
step 4 "构建项目..."
pnpm build
if [ $? -ne 0 ]; then
    error "项目构建失败。请检查构建日志获取详细信息"
fi
success "项目构建完成"

# 步骤 5: 创建系统服务
step 5 "创建系统服务..."
sudo tee /etc/systemd/system/servon.service << EOF
[Unit]
Description=Servon Server Panel
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$(pwd)
ExecStart=$(which node) .output/server/index.mjs
Restart=always

[Install]
WantedBy=multi-user.target
EOF

if [ $? -ne 0 ]; then
    error "无法创建服务文件。请检查是否有足够的权限"
fi
success "系统服务创建完成"

# 步骤 6: 启动服务
step 6 "启动服务..."
printf "  • 重新加载 systemd 配置...\n"
sudo systemctl daemon-reload
if [ $? -ne 0 ]; then
    error "无法重新加载 systemd 配置。请检查 systemd 状态"
fi

printf "  • 启用服务...\n"
sudo systemctl enable servon
if [ $? -ne 0 ]; then
    error "无法启用服务。请检查服务配置是否正确"
fi

printf "  • 启动服务...\n"
sudo systemctl start servon
if [ $? -ne 0 ]; then
    error "无法启动服务。请使用 'journalctl -u servon' 查看详细错误信息"
fi
success "服务启动完成"

# 安装完成
printf "\n"
printf "${GREEN}✨ Servon 安装成功!${NC}\n"
printf "你现在可以通过浏览器访问: ${BLUE}http://localhost:3000${NC}\n"
printf "如需查看服务状态，请运行: ${BLUE}sudo systemctl status servon${NC}\n"
printf "如需查看服务日志，请运行: ${BLUE}sudo journalctl -u servon -f${NC}\n" 