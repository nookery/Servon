<script setup lang="ts">
import { RiInboxLine } from '@remixicon/vue'
import { ref, provide } from 'vue'

interface Tab {
    key: string
    title: string
    icon?: any
}

const props = withDefaults(defineProps<{
    title?: string
    error?: string | null
    empty?: boolean
    emptyText?: string
    emptyDescription?: string
    emptyIcon?: any
    tabs?: Tab[]
    modelValue?: string
}>(), {
    modelValue: ''
})

const emit = defineEmits<{
    'update:modelValue': [value: string]
}>()

const activeTab = ref(props.modelValue || (props.tabs?.[0]?.key ?? ''))

// 提供给子组件的上下文
provide('activeTab', activeTab)

// 处理标签切换
function handleTabChange(tabKey: string) {
    activeTab.value = tabKey
    emit('update:modelValue', tabKey)
}
</script>

<template>
    <div class="card bg-base-100 h-full flex flex-col overflow-hidden">
        <!-- 固定头部 - 保持固定并添加阴影效果 -->
        <div class="flex-none bg-base-100 z-10 p-2 pb-0 shadow-sm">
            <!-- 标签页 -->
            <div v-if="tabs?.length" role="tablist" class="tabs tabs-lift bg-base-200 p-1 mb-4">
                <a role="tab" v-for="tab in tabs" :key="tab.key" class="tab gap-2"
                    :class="{ 'tab-active': activeTab === tab.key }" @click="handleTabChange(tab.key)">
                    <component v-if="tab.icon" :is="tab.icon" class="text-lg" />
                    {{ tab.title }}
                </a>
            </div>

            <!-- 头部内容插槽 -->
            <slot name="header"></slot>
        </div>

        <!-- 内容区域 - 移除所有滚动控制，只保持基础布局 -->
        <div class="flex-1 min-h-0 p-2 pt-4">
            <!-- 空状态显示 -->
            <div v-if="empty" class="card bg-base-200 p-8 text-center">
                <div class="flex flex-col items-center">
                    <component :is="emptyIcon || RiInboxLine" class="w-24 h-24 text-base-content/30 mb-4" />
                    <div class="text-xl mb-2">{{ emptyText || '暂无数据' }}</div>
                    <p v-if="emptyDescription" class="text-base-content/70">{{ emptyDescription }}</p>
                </div>
            </div>

            <!-- 内容区域 -->
            <div class="h-full">
                <!-- 默认内容 -->
                <template v-if="!tabs?.length">
                    <slot></slot>
                </template>

                <!-- 具名插槽用于不同标签页内容 -->
                <template v-else v-for="tab in tabs" :key="tab.key">
                    <div v-show="activeTab === tab.key" class="h-full">
                        <slot :name="tab.key"></slot>
                    </div>
                </template>
            </div>
        </div>
    </div>
</template>