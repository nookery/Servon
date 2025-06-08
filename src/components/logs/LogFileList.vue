<script setup lang="ts">
import { defineProps, defineEmits } from 'vue'
import type { LogFile } from '../../types/log'

defineProps<{
    logFiles: LogFile[]
    selectedFile: string
}>()

const emit = defineEmits<{
    'update:selectedFile': [value: string]
    'refresh': []
}>()

function selectFile(filePath: string) {
    emit('update:selectedFile', filePath)
}
</script>

<template>
    <div class="card bg-base-200 overflow-hidden h-full">
        <div class="h-full flex flex-col">
            <!-- 文件列表头部 -->
            <div class="flex-none p-3 flex justify-between items-center">
                <span class="text-sm font-medium">日志文件列表</span>
                <button @click="emit('refresh')" class="btn btn-ghost btn-xs">
                    <i class="ri-refresh-line"></i>
                </button>
            </div>

            <!-- 可滚动的文件列表 -->
            <div class="flex-1 min-h-0 overflow-y-auto p-2">
                <ul class="menu bg-base-200 w-full">
                    <li v-for="file in logFiles" :key="file.path">
                        <a class="flex items-center gap-2 transition-colors duration-200 hover:bg-base-300" :class="{
                            'bg-primary/10 text-primary border-l-4 border-primary': selectedFile === file.path,
                            'border-l-4 border-transparent': selectedFile !== file.path
                        }" @click="selectFile(file.path)">
                            <i class="ri-file-text-line" />
                            <span class="truncate">{{ file.path }}</span>
                        </a>
                    </li>
                    <li v-if="logFiles.length === 0" class="p-4 text-center text-base-content/50">
                        暂无日志文件
                    </li>
                </ul>
            </div>
        </div>
    </div>
</template>