.PHONY: build generate

# 默认生成所有内容
generate:
	@if [ "$(SKIP_GENERATE)" = "" ]; then \
		echo "🔨 Generating assets..." ; \
		go generate ./... ; \
	fi

# 构建整个项目
# LDFLAGS 可以从命令行传入，用于注入版本信息
build: generate
	go build -ldflags "$(LDFLAGS)" -o temp/servon

# 供 air 使用
air: 
	SKIP_GENERATE=1 make build

# 启动 Web 服务
serve: 
	go run main.go serve