FROM golang:1.23.6

# 安装 Node.js 和 npm
RUN apt-get update && apt-get install -y \
    nodejs npm zsh curl git

RUN npm install -g pnpm

# 安装 air
RUN go install github.com/air-verse/air@latest

# 安装 Oh My Zsh
RUN sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"

# 设置 zsh 为默认 shell
ENV SHELL /bin/zsh

ENTRYPOINT ["tail", "-f", "/dev/null"]
