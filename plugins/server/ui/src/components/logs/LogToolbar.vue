<script setup lang="ts">
import { defineProps, defineEmits, ref } from 'vue'
import IconButton from '../IconButton.vue'

defineProps<{
    currentDir: string
    selectedFile: string | null
}>()

const emit = defineEmits<{
    'update:currentDir': [value: string]
    'refresh': []
    'search': [keyword: string]
    'delete-log': []
    'clean-logs': []
    'clear-log': []
}>()

const searchKeyword = ref('')

function handleSearch() {
    emit('search', searchKeyword.value)
}

function updateDir(e: Event) {
    const target = e.target as HTMLInputElement
    emit('update:currentDir', target.value)
}
</script>

<template>
    <div class="flex flex-wrap justify-between items-center gap-4">
        <div class="flex gap-2 items-center flex-grow">
            <input type="text" :value="currentDir" @input="updateDir" placeholder="日志目录"
                class="input input-bordered input-sm flex-grow max-w-xs" />
            <IconButton icon="ri-refresh-line" size="sm" @click="emit('refresh')">刷新</IconButton>
        </div>
        <div class="flex flex-wrap gap-2">
            <div class="join">
                <input type="text" v-model="searchKeyword" placeholder="搜索日志"
                    class="input input-bordered input-sm join-item w-40 sm:w-auto" @keyup.enter="handleSearch" />
                <IconButton icon="ri-search-line" size="sm" class="join-item" @click="handleSearch">搜索</IconButton>
            </div>
            <div class="flex flex-wrap gap-2">
                <IconButton icon="ri-delete-bin-line" variant="error" size="sm" @click="emit('delete-log')"
                    :disabled="!selectedFile" title="删除当前日志文件">
                    <span class="hidden sm:inline">删除日志</span>
                    <span class="sm:hidden">删除</span>
                </IconButton>
                <IconButton icon="ri-delete-bin-line" variant="error" size="sm" @click="emit('clean-logs')"
                    title="清理30天前的日志">
                    <span class="hidden sm:inline">清理旧日志</span>
                    <span class="sm:hidden">清理</span>
                </IconButton>
                <IconButton icon="ri-eraser-line" variant="error" size="sm" @click="emit('clear-log')"
                    :disabled="!selectedFile" title="清空当前日志内容">
                    <span class="hidden sm:inline">清空日志</span>
                    <span class="sm:hidden">清空</span>
                </IconButton>
            </div>
        </div>
    </div>
</template>