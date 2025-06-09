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
	make build-frontend
	make build-backend

# 构建后端
# LDFLAGS 可以从命令行传入，用于注入版本信息
build-backend:
	go build -ldflags "$(LDFLAGS)" -o temp/servon

# 构建前端
build-frontend:
	cd plugins/server/ui  && pnpm install
	cd plugins/server/ui && pnpm build

dev:
	make build-backend
	cd plugins/server/ui && npm run dev

# 供 air 使用
air: 
	SKIP_GENERATE=1 make build-backend