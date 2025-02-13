FROM ubuntu:24.04

RUN apt-get update && apt-get install -y \
    nodejs npm zsh curl git golang

# 设置 Go 的工作目录
ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

# 创建工作目录
RUN mkdir -p $GOPATH/src $GOPATH/bin $GOPATH/pkg

RUN npm install -g pnpm

# 安装 air
RUN go install github.com/air-verse/air@latest

# 安装 Oh My Zsh
RUN sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"

# 设置 zsh 为默认 shell
ENV SHELL /bin/zsh

# 设置默认用户
USER root

ENTRYPOINT ["tail", "-f", "/dev/null"]
