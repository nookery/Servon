.PHONY: generate

# 生成步骤 - 执行 npm build 生成静态资源
generate:
	@if [ "$(SKIP_GENERATE)" = "" ]; then \
		pwd ; \
		pnpm install ; \
		pnpm build ; \
	fi