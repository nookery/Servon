<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'

const props = defineProps<{
    modelValue: 'table' | 'terminal'
}>()

const emit = defineEmits<{
    (e: 'update:modelValue', value: 'table' | 'terminal'): void
}>()

const viewMode = ref(props.modelValue)

// 监听 props 变化
watch(() => props.modelValue, (newValue) => {
    viewMode.value = newValue
})

// 监听内部状态变化并触发事件
watch(viewMode, (newValue) => {
    emit('update:modelValue', newValue)
    // 保存到本地存储
    localStorage.setItem('logViewMode', newValue)
})

// 组件挂载时从本地存储加载
onMounted(() => {
    const savedViewMode = localStorage.getItem('logViewMode') as 'table' | 'terminal' | null
    if (savedViewMode && (savedViewMode === 'table' || savedViewMode === 'terminal')) {
        viewMode.value = savedViewMode
        emit('update:modelValue', savedViewMode)
    }
})
</script>

<template>
    <div class="card bg-base-200 p-3">
        <h3 class="text-sm font-medium mb-2">视图模式</h3>
        <div class="grid grid-cols-2 gap-2">
            <button class="btn btn-sm flex items-center justify-center"
                :class="{ 'btn-primary': viewMode === 'table', 'btn-outline': viewMode !== 'table' }"
                @click="viewMode = 'table'">
                <i class="ri-table-line mr-1"></i>
                表格视图
            </button>
            <button class="btn btn-sm flex items-center justify-center"
                :class="{ 'btn-primary': viewMode === 'terminal', 'btn-outline': viewMode !== 'terminal' }"
                @click="viewMode = 'terminal'">
                <i class="ri-terminal-line mr-1"></i>
                终端视图
            </button>
        </div>
    </div>
</template>