#!/bin/bash

# 检查是否已安装 Node.js
if ! command -v node &> /dev/null; then
    echo "Installing Node.js..."
    curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
    sudo apt-get install -y nodejs
fi

# 检查是否已安装 pnpm
if ! command -v pnpm &> /dev/null; then
    echo "Installing pnpm..."
    curl -fsSL https://get.pnpm.io/install.sh | sh -
fi

# 克隆项目
echo "Cloning Servon..."
git clone https://github.com/your-username/servon.git
cd servon

# 安装依赖
echo "Installing dependencies..."
pnpm install

# 构建项目
echo "Building project..."
pnpm build

# 创建系统服务
echo "Creating system service..."
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

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable servon
sudo systemctl start servon

echo "Servon has been installed successfully!"
echo "You can access it at: http://localhost:3000" 