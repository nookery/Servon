<script setup lang="ts">
import { ref } from 'vue'
import IconButton from './IconButton.vue'
import { RiErrorWarningLine } from '@remixicon/vue'
import { useToast } from '../composables/useToast'

const errors = ref<{ id: number; message: string }[]>([])
const toast = useToast()
let nextId = 1

// 提供给其他组件使用的方法
const showError = (message: string) => {
    const id = nextId++
    errors.value.push({ id, message })
}

const removeError = (id: number) => {
    errors.value = errors.value.filter(error => error.id !== id)
}

const copyError = async (message: string) => {
    try {
        await navigator.clipboard.writeText(message)
        toast.success('错误信息已复制到剪贴板')
    } catch (err) {
        toast.error('复制失败，请手动复制')
    }
}

// 暴露方法给其他组件使用
defineExpose({
    showError
})
</script>

<template>
    <div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2 max-w-md">
        <div v-for="error in errors" :key="error.id" class="alert alert-error shadow-lg animate-slide-in-right">
            <div class="flex-1 flex items-start gap-2">
                <RiErrorWarningLine class="w-6 h-6 flex-shrink-0" />
                <span class="whitespace-pre-wrap">{{ error.message }}</span>
            </div>
            <div class="flex gap-1">
                <IconButton icon="ri-file-copy-line" variant="ghost" circle size="sm" title="复制错误信息"
                    @click="copyError(error.message)" />
                <IconButton icon="ri-close-line" variant="ghost" circle size="sm" title="关闭"
                    @click="removeError(error.id)" />
            </div>
        </div>
    </div>
</template>

<style scoped>
.animate-slide-in-right {
    animation: slide-in-right 0.3s ease-out;
}

@keyframes slide-in-right {
    from {
        transform: translateX(100%);
        opacity: 0;
    }

    to {
        transform: translateX(0);
        opacity: 1;
    }
}
</style>