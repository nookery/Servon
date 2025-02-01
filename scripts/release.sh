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

# 检查是否在 git 仓库中
if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
    print_error "Not in a git repository"
    exit 1
fi

# 检查工作目录是否干净
if ! git diff --quiet HEAD; then
    print_error "Working directory is not clean. Please commit or stash changes first."
    exit 1
fi

# 获取当前版本
CURRENT_VERSION=$(grep 'fmt.Println("Servon v' main.go | sed -E 's/.*Servon v([0-9]+\.[0-9]+\.[0-9]+).*/\1/')
if [ -z "$CURRENT_VERSION" ]; then
    print_error "Could not determine current version"
    exit 1
fi

print_info "Current version: v${CURRENT_VERSION}"

# 提示输入新版本
read -p "Enter new version (without v prefix): " NEW_VERSION

# 验证版本格式
if ! echo "$NEW_VERSION" | grep -E '^[0-9]+\.[0-9]+\.[0-9]+$' > /dev/null; then
    print_error "Invalid version format. Please use semantic versioning (e.g., 1.0.0)"
    exit 1
fi

# 更新版本号
print_info "Updating version to v${NEW_VERSION}..."

# 更新 main.go 中的版本号
sed -i.bak "s/Servon v[0-9]\+\.[0-9]\+\.[0-9]\+/Servon v${NEW_VERSION}/" main.go
rm -f main.go.bak

# 提交更改
git add main.go
git commit -m "chore: bump version to v${NEW_VERSION}"

# 创建 tag
git tag -a "v${NEW_VERSION}" -m "Release v${NEW_VERSION}"

print_success "Version updated to v${NEW_VERSION}"
print_info "To complete the release, run:"
echo
echo "    git push origin main v${NEW_VERSION}"
echo
print_warning "This will trigger the GitHub Actions workflow to create a new release" 