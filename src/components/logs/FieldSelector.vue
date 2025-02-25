<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useToast } from '../../composables/useToast'

const toast = useToast()

const props = defineProps<{
    modelValue: string[]
}>()

const emit = defineEmits<{
    (e: 'update:modelValue', value: string[]): void
}>()

const visibleFields = ref<string[]>([...props.modelValue])

// 监听 props 变化
watch(() => props.modelValue, (newValue) => {
    visibleFields.value = [...newValue]
})

// 处理字段显示切换，确保至少有一个字段可见
function toggleField(field: string) {
    // 如果要取消选中当前字段
    if (visibleFields.value.includes(field)) {
        // 检查是否只剩下这一个字段
        if (visibleFields.value.length === 1) {
            // 如果只剩一个字段，不允许取消选中
            toast.warning('至少需要显示一个字段')
            return
        }
        // 否则移除该字段
        visibleFields.value = visibleFields.value.filter(f => f !== field)
    } else {
        // 添加字段
        visibleFields.value.push(field)
    }

    // 更新父组件值
    emit('update:modelValue', [...visibleFields.value])

    // 保存到本地存储
    localStorage.setItem('logVisibleFields', JSON.stringify(visibleFields.value))
}

// 组件挂载时从本地存储加载
onMounted(() => {
    try {
        const savedFields = localStorage.getItem('logVisibleFields')
        if (savedFields) {
            const parsedFields = JSON.parse(savedFields) as string[]
            // 确保至少有一个字段
            if (parsedFields.length > 0) {
                visibleFields.value = parsedFields
                emit('update:modelValue', parsedFields)
            }
        }
    } catch (err) {
        console.error('Failed to load field preferences from localStorage', err)
    }
})
</script>

<template>
    <div class="card bg-base-200 p-3">
        <h3 class="text-sm font-medium mb-2">显示字段</h3>
        <div class="grid grid-cols-2 gap-2">
            <button class="flex items-center justify-center btn btn-sm"
                :class="{ 'btn-primary': visibleFields.includes('time'), 'btn-outline': !visibleFields.includes('time') }"
                @click="toggleField('time')">
                <i class="ri-time-line mr-1"></i>
                时间
            </button>
            <button class="flex items-center justify-center btn btn-sm"
                :class="{ 'btn-primary': visibleFields.includes('level'), 'btn-outline': !visibleFields.includes('level') }"
                @click="toggleField('level')">
                <i class="ri-dashboard-line mr-1"></i>
                级别
            </button>
            <button class="flex items-center justify-center btn btn-sm"
                :class="{ 'btn-primary': visibleFields.includes('caller'), 'btn-outline': !visibleFields.includes('caller') }"
                @click="toggleField('caller')">
                <i class="ri-code-line mr-1"></i>
                调用位置
            </button>
            <button class="flex items-center justify-center btn btn-sm"
                :class="{ 'btn-primary': visibleFields.includes('message'), 'btn-outline': !visibleFields.includes('message') }"
                @click="toggleField('message')">
                <i class="ri-message-2-line mr-1"></i>
                消息内容
            </button>
        </div>
    </div>
</template>