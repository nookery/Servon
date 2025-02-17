.PHONY: build generate

# 默认生成所有内容
generate:
	@echo "Generating assets..."
	@cd plugins/web && make generate
	@# 如果有其他需要生成的插件，在这里添加

# 构建整个项目
build: generate
	go build -o temp/servon

# 供 air 使用
air: 
	SKIP_GENERATE=1 make build

# 启动 Web 服务
serve: 
	go run main.go serve