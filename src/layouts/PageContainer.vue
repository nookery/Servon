<script setup lang="ts">
import { RiInboxLine, RiPushpinLine, RiPushpinFill } from '@remixicon/vue'
import { ref, provide, watch, onMounted, computed } from 'vue'

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
    // 添加一个唯一标识符，用于存储标签页状态
    layoutId?: string
}>(), {
    modelValue: '',
    layoutId: 'default'
})

const emit = defineEmits<{
    'update:modelValue': [value: string]
}>()

const activeTab = ref(props.modelValue || (props.tabs?.[0]?.key ?? ''))
const pinnedTabs = ref<Set<string>>(new Set()) // 存储被固定的标签页

// 计算当前显示的标签页
const visibleTabs = computed(() => {
    const result = new Set(pinnedTabs.value)
    if (!result.has(activeTab.value)) {
        result.add(activeTab.value)
    }
    return result
})

// 计算布局列数
const gridColumns = computed(() => {
    const size = visibleTabs.value.size
    if (size <= 1) return 1
    if (size <= 4) return 2
    return Math.ceil(Math.sqrt(size)) // 动态计算列数
})

// 提供给子组件的上下文
provide('activeTab', activeTab)
provide('pinnedTabs', pinnedTabs)

// 从 localStorage 加载标签页状态
onMounted(() => {
    // 只有当提供了 layoutId 时才尝试恢复布局
    if (props.layoutId) {
        const savedLayout = localStorage.getItem(`layout_${props.layoutId}`)
        if (savedLayout) {
            try {
                const layoutData = JSON.parse(savedLayout)

                // 恢复固定的标签页
                if (layoutData.pinnedTabs && Array.isArray(layoutData.pinnedTabs)) {
                    pinnedTabs.value = new Set(layoutData.pinnedTabs)
                }

                // 恢复活动标签页
                if (layoutData.activeTab && props.tabs?.some(tab => tab.key === layoutData.activeTab)) {
                    activeTab.value = layoutData.activeTab
                    emit('update:modelValue', activeTab.value)
                }
            } catch (e) {
                console.error('Error restoring layout:', e)
            }
        }
    }
})

// 监听标签页状态变化并保存
watch([activeTab, pinnedTabs], () => {
    if (props.layoutId) {
        const layoutData = {
            activeTab: activeTab.value,
            pinnedTabs: Array.from(pinnedTabs.value)
        }
        localStorage.setItem(`layout_${props.layoutId}`, JSON.stringify(layoutData))
    }
}, { deep: true })

// 监听modelValue变化
watch(() => props.modelValue, (newVal) => {
    if (newVal && newVal !== activeTab.value) {
        activeTab.value = newVal
    }
})

// 处理标签切换
function handleTabChange(tabKey: string) {
    activeTab.value = tabKey
    emit('update:modelValue', tabKey)
}

// 切换标签页固定状态
function togglePinned(tabKey: string) {
    const newPinnedTabs = new Set(pinnedTabs.value)

    if (newPinnedTabs.has(tabKey)) {
        newPinnedTabs.delete(tabKey)
    } else {
        newPinnedTabs.add(tabKey)
    }

    pinnedTabs.value = newPinnedTabs
    console.log('标签页固定状态切换为:', tabKey, newPinnedTabs.has(tabKey) ? '已固定' : '未固定')
}

// 检查标签页是否被固定
function isTabPinned(tabKey: string) {
    return pinnedTabs.value.has(tabKey)
}

// 获取标签对象
function getTabByKey(key: string) {
    return props.tabs?.find(tab => tab.key === key)
}
</script>

<template>
    <div class="card bg-base-100 h-full flex flex-col overflow-hidden">
        <!-- 固定头部 - 保持固定并添加阴影效果 -->
        <div class="flex-none bg-base-100 z-10 p-2 pb-0 shadow-sm">
            <!-- 标签页 -->
            <div v-if="tabs?.length" class="flex justify-between items-center mb-4">
                <div role="tablist" class="tabs tabs-lift bg-base-200 p-1">
                    <a role="tab" v-for="tab in tabs" :key="tab.key" class="tab gap-2" :class="{
                        'tab-active': activeTab === tab.key,
                        'font-bold': isTabPinned(tab.key)
                    }" @click="handleTabChange(tab.key)">
                        <component v-if="tab.icon" :is="tab.icon" class="text-lg" />
                        {{ tab.title }}
                        <RiPushpinFill v-if="isTabPinned(tab.key)" class="text-xs text-primary" />
                    </a>
                </div>
            </div>

            <!-- 头部内容插槽 -->
            <slot name="header"></slot>
        </div>

        <!-- 内容区域 -->
        <div class="flex-1 min-h-0 p-2 pt-4 overflow-auto">
            <!-- 空状态显示 -->
            <div v-if="empty" class="card bg-base-200 p-8 text-center">
                <div class="flex flex-col items-center">
                    <component :is="emptyIcon || RiInboxLine" class="w-24 h-24 text-base-content/30 mb-4" />
                    <div class="text-xl mb-2">{{ emptyText || '暂无数据' }}</div>
                    <p v-if="emptyDescription" class="text-base-content/70">{{ emptyDescription }}</p>
                </div>
            </div>

            <!-- 内容区域 -->
            <div v-if="!empty" class="h-full grid gap-4" :style="{
                'grid-template-columns': `repeat(${gridColumns}, minmax(0, 1fr))`
            }">
                <!-- 默认内容 -->
                <template v-if="!tabs?.length">
                    <slot></slot>
                </template>

                <!-- 具名插槽用于不同标签页内容 -->
                <template v-else>
                    <!-- 可见的标签页 -->
                    <template v-for="tabKey in Array.from(visibleTabs)" :key="tabKey">
                        <div class="h-full flex flex-col relative">
                            <div class="absolute top-2 right-2 z-10">
                                <button @click="togglePinned(tabKey)" class="btn btn-xs btn-circle">
                                    <component :is="isTabPinned(tabKey) ? RiPushpinFill : RiPushpinLine" class="text-lg"
                                        :class="{ 'text-primary': isTabPinned(tabKey) }" />
                                </button>
                            </div>
                            <div class="card border rounded-lg p-4 h-full overflow-hidden relative" :class="{
                                'border-primary': isTabPinned(tabKey),
                                'bg-base-200': activeTab === tabKey && !isTabPinned(tabKey)
                            }">
                                <div class="text-lg font-medium mb-2 pb-2 border-b flex items-center gap-2">
                                    <component v-if="getTabByKey(tabKey)?.icon" :is="getTabByKey(tabKey)?.icon"
                                        class="text-lg" />
                                    {{ getTabByKey(tabKey)?.title }}
                                </div>
                                <div class="flex-1 overflow-auto">
                                    <slot :name="tabKey"></slot>
                                </div>
                            </div>
                        </div>
                    </template>
                </template>
            </div>
        </div>
    </div>
</template>