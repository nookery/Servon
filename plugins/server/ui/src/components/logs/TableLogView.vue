<script setup lang="ts">
import { RiEmotionHappyLine } from '@remixicon/vue'
import type { LogEntry } from '../../types/log'

defineProps<{
    logEntries: LogEntry[],
    visibleFields: string[]
}>()

// 格式化日志级别样式
function getLevelClass(level: string): string {
    const levelMap: Record<string, string> = {
        error: 'badge-error',
        warn: 'badge-warning',
        info: 'badge-info',
        debug: 'badge-neutral'
    }
    return `badge ${levelMap[level.toLowerCase()] || 'badge-ghost'}`
}
</script>

<template>
    <table class="table table-xs w-full">
        <thead class="sticky top-0 bg-base-200 z-10">
            <tr>
                <th v-if="visibleFields.includes('time')">时间</th>
                <th v-if="visibleFields.includes('level')">级别</th>
                <th v-if="visibleFields.includes('caller')">调用位置</th>
                <th v-if="visibleFields.includes('message')">消息</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="entry in logEntries" :key="entry.time">
                <td v-if="visibleFields.includes('time')" class="whitespace-nowrap">{{ entry.time }}</td>
                <td v-if="visibleFields.includes('level')">
                    <span :class="getLevelClass(entry.level)">
                        {{ entry.level }}
                    </span>
                </td>
                <td v-if="visibleFields.includes('caller')" class="text-xs">{{ entry.caller }}</td>
                <td v-if="visibleFields.includes('message')" class="whitespace-pre-wrap">{{ entry.message }}</td>
            </tr>
            <tr v-if="logEntries.length === 0">
                <td :colspan="visibleFields.length" class="text-center text-base-content/50 py-8">
                    <div class="flex flex-col items-center gap-2">
                        <RiEmotionHappyLine class="w-8 h-8 text-success" />
                        <span>暂无符合条件的日志 </span>
                        <span class="text-xs opacity-50">这说明系统运行得很顺利呢！</span>
                    </div>
                </td>
            </tr>
        </tbody>
    </table>
</template>