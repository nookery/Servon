<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'

const props = defineProps<{
    modelValue: string[]
}>()

const emit = defineEmits<{
    (e: 'update:modelValue', value: string[]): void
}>()

const selectedLevels = ref<string[]>([...props.modelValue])

// 监听 props 变化
watch(() => props.modelValue, (newValue) => {
    selectedLevels.value = [...newValue]
})

// 监听内部状态变化并触发事件
watch(selectedLevels, (newValue) => {
    emit('update:modelValue', [...newValue])
    // 保存到本地存储
    localStorage.setItem('logSelectedLevels', JSON.stringify(newValue))
})

// 组件挂载时从本地存储加载
onMounted(() => {
    try {
        const savedLevels = localStorage.getItem('logSelectedLevels')
        if (savedLevels) {
            const parsedLevels = JSON.parse(savedLevels) as string[]
            // 确保至少选择了一个级别
            if (parsedLevels.length > 0) {
                selectedLevels.value = parsedLevels
                emit('update:modelValue', parsedLevels)
            }
        }
    } catch (err) {
        console.error('Failed to load level preferences from localStorage', err)
    }
})
</script>

<template>
    <div class="card bg-base-200 p-3">
        <h3 class="text-sm font-medium mb-2">日志级别筛选</h3>
        <div class="grid grid-cols-2 gap-2">
            <label class="flex items-center justify-center btn btn-sm"
                :class="{ 'btn-error': selectedLevels.includes('error') }">
                <input type="checkbox" class="hidden" v-model="selectedLevels" value="error" />
                <i class="ri-error-warning-fill mr-1"></i>
                错误
            </label>
            <label class="flex items-center justify-center btn btn-sm"
                :class="{ 'btn-warning': selectedLevels.includes('warn') }">
                <input type="checkbox" class="hidden" v-model="selectedLevels" value="warn" />
                <i class="ri-alert-fill mr-1"></i>
                警告
            </label>
            <label class="flex items-center justify-center btn btn-sm"
                :class="{ 'btn-info': selectedLevels.includes('info') }">
                <input type="checkbox" class="hidden" v-model="selectedLevels" value="info" />
                <i class="ri-information-fill mr-1"></i>
                信息
            </label>
            <label class="flex items-center justify-center btn btn-sm"
                :class="{ 'btn-neutral': selectedLevels.includes('debug') }">
                <input type="checkbox" class="hidden" v-model="selectedLevels" value="debug" />
                <i class="ri-bug-fill mr-1"></i>
                调试
            </label>
        </div>
    </div>
</template>