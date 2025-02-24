<script setup lang="ts">
import Alert from '../components/Alert.vue'
import { RiInboxLine } from '@remixicon/vue'

defineProps<{
    title?: string
    error?: string | null
    empty?: boolean
    emptyText?: string
    emptyDescription?: string
    emptyIcon?: any // 允许传入自定义图标组件
}>()
</script>

<template>
    <div class="card bg-base-100 p-0 h-full flex flex-col">
        <div class="card-body p-0 flex flex-col h-full">
            <!-- 固定头部 -->
            <div class="sticky top-0 bg-base-100 z-10 p-2 pb-0">
                <!-- 错误提示 -->
                <Alert v-if="error" type="error" :message="error" class="mb-4" />

                <!-- 头部内容插槽 -->
                <slot name="header"></slot>
                <!-- Tab栏插槽 -->
                <slot name="tabs"></slot>
            </div>

            <!-- 可滚动的内容区域 -->
            <div class="flex-1 overflow-auto p-2 pt-0">
                <!-- 空状态显示 -->
                <div v-if="empty" class="card bg-base-200 p-8 text-center">
                    <div class="flex flex-col items-center">
                        <component :is="emptyIcon || RiInboxLine" class="w-24 h-24 text-base-content/30 mb-4" />
                        <div class="text-xl mb-2">{{ emptyText || '暂无数据' }}</div>
                        <p v-if="emptyDescription" class="text-base-content/70">{{ emptyDescription }}</p>
                    </div>
                </div>
                <!-- 默认内容 -->
                <template v-else>
                    <slot></slot>
                </template>
            </div>
        </div>
    </div>
</template>