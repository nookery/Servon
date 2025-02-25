<script setup lang="ts">
import { RiInboxLine, RiLayoutGridLine, RiLayoutLine } from '@remixicon/vue'
import { ref, provide, watch, onMounted } from 'vue'

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
    // 添加一个唯一标识符，用于存储布局模式
    layoutId?: string
}>(), {
    modelValue: '',
    layoutId: 'default'
})

const emit = defineEmits<{
    'update:modelValue': [value: string]
}>()

const activeTab = ref(props.modelValue || (props.tabs?.[0]?.key ?? ''))
const isGridView = ref(false)

// 提供给子组件的上下文
provide('activeTab', activeTab)

// 从 localStorage 加载布局模式
onMounted(() => {
    // 只有当提供了 layoutId 时才尝试恢复布局
    if (props.layoutId) {
        const savedLayout = localStorage.getItem(`layout_${props.layoutId}`)
        if (savedLayout) {
            try {
                const layoutData = JSON.parse(savedLayout)
                isGridView.value = layoutData.isGridView || false

                // 如果有保存的标签页，且该标签页在当前可用的标签页中存在，则恢复
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

// 监听布局变化并保存
watch([isGridView, activeTab], () => {
    if (props.layoutId) {
        const layoutData = {
            isGridView: isGridView.value,
            activeTab: activeTab.value
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

// 切换视图模式
function toggleViewMode() {
    isGridView.value = !isGridView.value
    console.log('视图模式切换为:', isGridView.value ? '网格视图' : '列表视图')
}
</script>

<template>
    <div class="card bg-base-100 h-full flex flex-col overflow-hidden">
        <!-- 固定头部 - 保持固定并添加阴影效果 -->
        <div class="flex-none bg-base-100 z-10 p-2 pb-0 shadow-sm">
            <!-- 标签页和视图切换按钮 -->
            <div v-if="tabs?.length" class="flex justify-between items-center mb-4">
                <div role="tablist" class="tabs tabs-lift bg-base-200 p-1">
                    <a role="tab" v-for="tab in tabs" :key="tab.key" class="tab gap-2"
                        :class="{ 'tab-active': activeTab === tab.key }" @click="handleTabChange(tab.key)">
                        <component v-if="tab.icon" :is="tab.icon" class="text-lg" />
                        {{ tab.title }}
                    </a>
                </div>

                <!-- 视图切换按钮 -->
                <button class="btn btn-sm" :class="isGridView ? 'btn-primary' : 'btn-outline'" @click="toggleViewMode"
                    title="切换视图模式">
                    <component :is="isGridView ? RiLayoutLine : RiLayoutGridLine" class="text-lg mr-1" />
                    {{ isGridView ? '列表视图' : '网格视图' }}
                </button>
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
            <div v-if="!empty" class="h-full" :class="{
                'grid grid-cols-1 md:grid-cols-2 gap-4': isGridView
            }">
                <!-- 默认内容 -->
                <template v-if="!tabs?.length">
                    <slot></slot>
                </template>

                <!-- 具名插槽用于不同标签页内容 -->
                <template v-else>
                    <template v-for="tab in tabs" :key="tab.key">
                        <div v-show="isGridView || activeTab === tab.key" class="h-full flex flex-col"
                            :class="{ 'border rounded-lg p-4 border-primary': isGridView }">
                            <div v-if="isGridView" class="flex-none text-lg font-medium mb-2 pb-2 border-b">{{ tab.title
                            }}</div>
                            <div class="flex-1 overflow-auto">
                                <slot :name="tab.key"></slot>
                            </div>
                        </div>
                    </template>
                </template>
            </div>
        </div>
    </div>
</template>