# 构建系统

## Makefile简介

Makefile是一个用于自动化构建过程的工具，它定义了项目中的各种构建任务及其依赖关系。Make工具会读取Makefile文件，并根据定义的规则执行相应的命令。

### 为什么使用Makefile

1. **自动化构建** - 简化复杂的构建流程
2. **增量构建** - 只重新构建发生变化的部分
3. **标准化流程** - 确保所有开发者使用相同的构建步骤
4. **可维护性** - 集中管理构建逻辑

### 基本语法和特殊标记

Makefile主要由以下部分组成：
- **目标(Target)** - 定义要执行的任务名称（如 build, generate）
- **依赖(Prerequisites)** - 指定目标依赖的其他目标或文件
- **命令(Recipe)** - 实现目标的shell命令

一些重要的语法特性：

1. **.PHONY标记** - 用于声明虚拟目标，表示该目标不是实际的文件名
   ```makefile
   .PHONY: build generate
   ```

2. **@符号** - 命令前加@表示执行时不显示该命令本身
   ```makefile
   target:
       @echo "构建中..." # 只显示"构建中..."，不显示echo命令本身
   ```

3. **变量使用** - 使用$(变量名)引用变量
   ```makefile
   LDFLAGS=-X main.Version=1.0.0
   build:
       go build -ldflags "$(LDFLAGS)"
   ```

4. **条件判断** - 使用shell的if语句进行条件判断
   ```makefile
   target:
       @if [ "$(SKIP_GENERATE)" = "" ]; then \
           echo "执行生成步骤" ; \
       fi
   ```

## 项目任务说明

本项目使用Make作为构建工具。Makefile中定义了以下主要任务：

## 主要任务

### `make generate`

生成项目所需的资源文件。如果不希望执行生成步骤，可以设置环境变量 `SKIP_GENERATE=1`。

### `make build`

构建整个项目，包括前端和后端。这个任务会：
1. 执行 generate 任务生成资源
2. 构建前端代码
3. 构建后端代码

### `make build-backend`

仅构建后端部分。支持通过 `LDFLAGS` 参数注入版本信息，例如：
```bash
make build-backend LDFLAGS="-X main.Version=1.0.0"
```

### `make build-frontend`

仅构建前端部分。这个任务会：
1. 安装前端依赖（使用pnpm）
2. 构建前端代码

### `make air`

用于开发时的快速构建。这个任务会：
1. 跳过资源生成步骤（设置 SKIP_GENERATE=1）
2. 仅构建后端代码

这个任务主要配合 air 工具使用，适合在开发过程中进行快速的代码修改和测试。